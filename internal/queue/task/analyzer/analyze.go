package analyzer

import (
	"context"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/semgroup"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/h2non/filetype"

	"github.com/labstack/gommon/log"
	"github.com/marktrs/gitsast/app"
	"github.com/marktrs/gitsast/internal/model"
	ahocorasick "github.com/petar-dambovaliev/aho-corasick"
	"github.com/vmihailenco/taskq/v3"
)

var _ IAnalyzeTask = (*Analyzer)(nil)

var cloneLocationPrefix = "temp/"

// Task - register analyze task into task queue
var (
	Task = taskq.RegisterTask(&taskq.TaskOptions{
		Name: "analyzer",
		Handler: func(app *app.App, id string) error {
			return NewAnalyzer(app).Analyze(id)
		},
	})
)

// IAnalyzeTask - interface for analyze task
type IAnalyzeTask interface {
	Analyze(reportId string) error
}

type Analyzer struct {
	repo   model.IRepositoryRepo
	report model.IReportRepo
	rule   model.IRuleRepo
}

// NewAnalyzeTask - create new analyze task
func NewAnalyzer(app *app.App) IAnalyzeTask {
	// connect db
	db := app.DB()
	repo := model.NewRepositoryRepo(db)
	report := model.NewReportRepo(db)
	rule := model.NewRuleRepo(db)

	return &Analyzer{repo, report, rule}
}

// Analyze - implement analyze task interface
func (t *Analyzer) Analyze(reportId string) error {
	ctx := context.Background()

	report, err := t.report.GetById(ctx, reportId)
	if err != nil {
		log.Error(err)
		return err
	}

	report.StartedAt = time.Now()

	// set status
	log.Infof("starting analyzed task report_id=%s", report.ID)
	if err := t.setReportStatus(report, model.StatusInProgress); err != nil {
		return t.handleFailedTask(report, err)
	}

	repo, err := t.repo.GetById(ctx, report.RepositoryID)
	if err != nil {
		return t.handleFailedTask(report, err)
	}

	// look up for latest rules
	rules, err := t.rule.GetAll(ctx)
	if err != nil {
		return t.handleFailedTask(report, err)
	}

	tmpDir := path.Join(cloneLocationPrefix, repo.ID)

	log.Infof("getting paths from remote url=%s", repo.RemoteURL)
	// look up for latest code
	paths, err := t.getPathsFromRemoteURL(tmpDir, repo.RemoteURL)
	if err != nil {
		return t.handleFailedTask(report, err)
	}

	log.Infof("scanning files for issues url=%s", repo.RemoteURL)
	issues, err := t.scanFilesForIssues(repo.ID, paths, rules)
	if err != nil {
		return t.handleFailedTask(report, err)
	}

	t.removeTempDir(tmpDir)

	log.Infof("adding issues to report report_id=%s", report.ID)

	report.Issues = issues
	report.FinishedAt = time.Now()

	if err := t.setReportStatus(report, model.StatusSuccess); err != nil {
		return t.handleFailedTask(report, err)
	}

	log.Infof("analyzed task completed report_id=%s", report.ID)

	return nil
}

func (t *Analyzer) removeTempDir(tmpDir string) {
	if err := os.RemoveAll(tmpDir); err != nil {
		log.Error(err)
	}
}

// handleFailedTask - set report status to failed with reason
func (t *Analyzer) handleFailedTask(report *model.Report, err error) error {
	log.Error(err)
	report.FailedReason = err.Error()
	return t.setReportStatus(report, model.StatusSuccess)
}

// setReportStatus - set report status
func (t *Analyzer) setReportStatus(report *model.Report, status model.ReportStatus) error {
	ctx := context.Background()

	now := time.Now()

	report.Status = status
	report.UpdatedAt = now

	switch status {
	case model.StatusInProgress:
		report.StartedAt = now
	case model.StatusSuccess:
		report.FinishedAt = now
	case model.StatusFailed:
		report.FinishedAt = now
	}

	_, err := t.report.Update(ctx, report)
	if err != nil {
		return err
	}

	return nil
}

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

func (t *Analyzer) detectIssueLocation(fragment Fragment, rule *model.Rule) []*model.Issue {
	issues := make([]*model.Issue, 0)
	regex := regexp.MustCompile(rule.Keyword)
	matchIndices := regex.FindAllStringIndex(fragment.Raw, -1)

	for _, matchIndex := range matchIndices {
		loc := location(fragment, matchIndex)

		if matchIndex[1] > loc.endLineIndex {
			loc.endLineIndex = matchIndex[1]
		}

		issues = append(issues, &model.Issue{
			RuleID: model.GetFormattedRuleId(rule.ID),
			Location: model.Location{
				Path: fragment.FilePath,
				Line: uint64(loc.startLine),
			},
			Description: rule.Description,
			Severity:    rule.Severity.String(),
			Keyword:     rule.Keyword,
		})
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
