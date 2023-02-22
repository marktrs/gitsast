package analyzer

import (
	"context"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/semgroup"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/h2non/filetype"
	"github.com/labstack/gommon/log"
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

// getPathsFromRemoteURL - get all file paths from remote url except ignored file types
func (t *Analyzer) getPathsFromRemoteURL(tmpDir string, remoteURL string) ([]string, error) {
	// local clone
	r, err := git.PlainClone(tmpDir, false, &git.CloneOptions{URL: remoteURL})
	if err != nil {
		return nil, err
	}

	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	treeWalker := object.NewTreeWalker(tree, true, nil)

	filepaths := make([]string, 0)
	for {
		name, _, err := treeWalker.Next()
		if err == io.EOF {
			break
		}

		isIgnored := isFileTypeIgnored(name)
		if isIgnored {
			continue
		}

		info, err := os.Stat(path.Join(tmpDir, name))
		if err != nil {
			log.Error(err)
			break
		}

		if info.IsDir() {
			continue
		}

		filepaths = append(filepaths, path.Join(tmpDir, name))
	}
	defer treeWalker.Close()

	return filepaths, nil
}

// scanFilesForIssues - scan files for issues
func (t *Analyzer) scanFilesForIssues(tmpDir string, paths []string, rules []*model.Rule) ([]*model.Issue, error) {
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

			for _, issue := range t.scanLineForIssue(fragment, rules) {
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
func (t *Analyzer) scanLineForIssue(fragment Fragment, rules []*model.Rule) []*model.Issue {
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
			issues = append(issues, t.detectIssueLocation(fragment, rule)...)
		}
	}

	return issues
}

var ignoredFileTypes = []string{
	".aac", ".aiff", ".ape", ".au", ".flac", ".gsm", ".it", ".m3u", ".m4a", ".mid", ".mod", ".mp3", ".mpa", ".pls", ".ra", ".s3m", ".sid", ".wav", ".wma", ".xm", ".7z",
	".a", ".ar", ".bz2", ".cab", ".cpio", ".deb", ".dmg", ".egg", ".gz", ".iso", ".lha", ".mar", ".pea", ".rar", ".rpm", ".s7z", ".shar", ".tar", ".tbz2", ".tgz", ".tlz",
	".whl", ".xpi", ".deb", ".rpm", ".xz", ".pak", ".crx", ".exe", ".msi", ".bin", ".eot", ".otf", ".ttf", ".woff", ".woff2", ".3dm", ".3ds", ".max", ".bmp", ".dds", ".gif",
	".jpg", ".jpeg", ".png", ".psd", ".xcf", ".tga", ".thm", ".tif", ".tiff", ".yuv", ".ai", ".eps", ".ps", ".svg", ".dwg", ".dxf", ".gpx", ".kml", ".kmz", ".ods", ".xls",
	".xlsx", ".csv", ".ics", ".vcf", ".ppt", ".odp", ".3g2", ".3gp", ".aaf", ".asf", ".avchd", ".avi", ".drc", ".flv", ".m2v", ".m4p", ".m4v", ".mkv", ".mng", ".mov", ".mp2",
	".mp4", ".mpe", ".mpeg", ".mpg", ".mpv", ".mxf", ".nsv", ".ogg", ".ogv", ".ogm", ".qt", ".rm", ".rmvb", ".roq", ".srt", ".svi", ".vob", ".webm", ".wmv", ".yuv",
}

func isFileTypeIgnored(filename string) bool {
	var isIgnored bool
	for _, ignoredFileType := range ignoredFileTypes {
		if strings.HasSuffix(filename, ignoredFileType) {
			isIgnored = true
		}
	}

	return isIgnored
}
