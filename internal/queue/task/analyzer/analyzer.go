package analyzer

import (
	"context"
	"os"
	"path"
	"time"

	"github.com/marktrs/gitsast/app"
	"github.com/marktrs/gitsast/internal/model"
	"github.com/marktrs/gitsast/internal/queue/task/analyzer/git"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/taskq/v3"
)

var _ IAnalyzeTask = (*Analyzer)(nil)

var cloneLocationPrefix = "temp/"

// Task - register analyze task into task queue
var (
	Task = taskq.RegisterTask(&taskq.TaskOptions{
		Name: "analyzer",
		Handler: func(reportId string) error {
			_, app, err := app.Start(context.Background(), "analyzerTask", "")
			if err != nil {
				return err
			}

			repo := model.NewRepositoryRepo(app)
			report := model.NewReportRepo(app)
			rule := model.NewRuleRepo(app)
			git := git.NewClient()
			detector := NewDetector()
			scanner := NewScanner(detector)

			a, err := NewAnalyzer(app, repo, report, rule, git, detector, scanner)
			if err != nil {
				return err
			}

			return a.Analyze(reportId)
		},
	})
)

// IAnalyzeTask - interface for analyze task
type IAnalyzeTask interface {
	Analyze(reportId string) error
}

type Analyzer struct {
	app      *app.App
	repo     model.IRepositoryRepo
	report   model.IReportRepo
	rule     model.IRuleRepo
	git      git.IClient
	detector Detector
	scanner  Scanner
}

// NewAnalyzeTask - create new analyze task
func NewAnalyzer(
	app *app.App,
	repo model.IRepositoryRepo,
	report model.IReportRepo,
	rule model.IRuleRepo,
	git git.IClient,
	detector Detector,
	scanner Scanner,

) (IAnalyzeTask, error) {

	return &Analyzer{app, repo, report, rule, git, detector, scanner}, nil
}

// Analyze - implement analyze task interface
func (a *Analyzer) Analyze(reportId string) error {
	ctx := context.Background()
	log := log.Info().Fields(map[string]interface{}{
		"report_id": reportId,
	})

	log.Msg("starting analyze task")

	report, err := a.report.GetById(ctx, reportId)
	if err != nil {
		log.Err(err)
		return err
	}

	report.StartedAt = time.Now()

	// set status
	log.Msg("set status to in_progress")
	if err := a.setReportStatus(report, model.StatusInProgress); err != nil {
		return a.handleFailedTask(report, err)
	}

	log.Str("repository_id", report.RepositoryID).Msg("get repository by id")
	repo, err := a.repo.GetById(ctx, report.RepositoryID)
	if err != nil {
		return a.handleFailedTask(report, err)
	}

	// look up for latest rules
	rules, err := a.rule.GetAll(ctx)
	if err != nil {
		return a.handleFailedTask(report, err)
	}

	tmpDir := path.Join(cloneLocationPrefix, repo.ID)

	log.Str("url", repo.RemoteURL).Msg("getting paths from remote url")
	paths, err := a.git.GetPathsFromRemoteURL(tmpDir, repo.RemoteURL)
	if err != nil {
		return a.handleFailedTask(report, err)
	}
	defer a.removeTempDir(tmpDir)

	log.Str("url", repo.RemoteURL).Msg("scanning files for issues")
	issues, err := a.scanner.ScanFilesForIssues(repo.ID, paths, rules)
	if err != nil {
		return a.handleFailedTask(report, err)
	}

	if len(issues) == 0 {
		log.Msg("no issues found")
	} else {
		log.Msg("adding issues to report")
		report.Issues = issues
	}

	report.FinishedAt = time.Now()
	if err := a.setReportStatus(report, model.StatusSuccess); err != nil {
		return a.handleFailedTask(report, err)
	}
	log.Msg("analyzed task completed")

	return nil
}

// removeTempDir - remove cloned repo directory
func (a *Analyzer) removeTempDir(tmpDir string) error {
	if err := os.RemoveAll(tmpDir); err != nil {
		log.Err(err)
		return err
	}

	return nil
}

// handleFailedTask - set report status to failed with reason
func (a *Analyzer) handleFailedTask(report *model.Report, err error) error {
	log.Err(err)
	report.FailedReason = err.Error()
	return a.setReportStatus(report, model.StatusFailed)
}

// setReportStatus - set report status
func (a *Analyzer) setReportStatus(report *model.Report, status model.ReportStatus) error {
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

	_, err := a.report.Update(ctx, report)
	if err != nil {
		return err
	}

	return nil
}
