package analyzer

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/semgroup"
	"github.com/h2non/filetype"
	"github.com/marktrs/gitsast/internal/model"
	ahocorasick "github.com/petar-dambovaliev/aho-corasick"
)

// Fragment represents a fragment of a file or a commit
type Fragment struct {
	// Raw is the raw content of the fragment
	Raw string
	// FilePath is the path to the file if applicable
	FilePath    string
	SymlinkFile string
	// CommitSHA is the SHA of the commit if applicable
	CommitSHA string
	// newlineIndices is a list of indices of newlines in the raw content.
	// This is used to calculate the line location of a finding
	newlineIndices [][]int
	// keywords is a map of all the keywords contain within the contents
	// of this fragment
	keywords map[string]bool
}

// Scanner represents a scanner
type Scanner interface {
	ScanFilesForIssues(tmpDir string, paths []string, rules []*model.Rule) ([]*model.Issue, error)
	ScanLineForIssues(fragment Fragment, rules []*model.Rule) []*model.Issue
}

type scanner struct {
	detector Detector
}

func NewScanner(detector Detector) Scanner {
	return &scanner{
		detector,
	}
}

// scanFilesForIssues - scan files for issues
func (sc *scanner) ScanFilesForIssues(tmpDir string, paths []string, rules []*model.Rule) ([]*model.Issue, error) {
	issues := make([]*model.Issue, 0)
	s := semgroup.NewGroup(context.Background(), 4)
	for _, path := range paths {
		s.Go(func() error {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			mimetype, err := filetype.Match(b)
			if err != nil {
				return err
			}

			if mimetype.MIME.Type == "application" {
				return nil // skip binary files
			}

			fragment := Fragment{
				Raw:      string(b),
				FilePath: path,
			}

			for _, issue := range sc.ScanLineForIssues(fragment, rules) {
				issue.Location.Path = strings.ReplaceAll(path, filepath.Join(cloneLocationPrefix, tmpDir), "")
				issues = append(issues, issue)
			}

			return nil
		})
	}

	if err := s.Wait(); err != nil {
		return issues, err
	}

	return issues, nil
}

// scanLineForIssue - scan a line for issues
func (sc *scanner) ScanLineForIssues(fragment Fragment, rules []*model.Rule) []*model.Issue {
	issues := make([]*model.Issue, 0)

	builder := ahocorasick.NewAhoCorasickBuilder(ahocorasick.Opts{
		AsciiCaseInsensitive: true,
		MatchOnlyWholeWords:  false,
		MatchKind:            ahocorasick.LeftMostLongestMatch,
		DFA:                  true,
	})

	keywords := make([]string, len(rules))
	for _, rule := range rules {
		keywords = append(keywords, rule.Keyword)
	}

	// build prefilter
	preFilter := builder.Build(keywords)
	// initiate fragment keywords
	fragment.keywords = make(map[string]bool)
	// add newline indices for location calculation in detectRule
	fragment.newlineIndices = regexp.MustCompile("\n").FindAllStringIndex(fragment.Raw, -1)

	// build keyword map for prefilter rules
	normalizedRaw := strings.ToLower(fragment.Raw)
	matches := preFilter.FindAll(normalizedRaw)

	for _, m := range matches {
		fragment.keywords[normalizedRaw[m.Start():m.End()]] = true
	}

	for _, rule := range rules {
		fragmentContainsKeyword := false
		// check if keyword are in the fragment
		if _, ok := fragment.keywords[strings.ToLower(rule.Keyword)]; ok {
			fragmentContainsKeyword = true
		}

		if fragmentContainsKeyword {
			issues = append(issues, sc.detector.DetectIssueLocation(fragment, rule)...)
		}
	}

	return issues
}
