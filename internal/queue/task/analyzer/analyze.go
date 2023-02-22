package analyzer

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/marktrs/gitsast/app"
	"github.com/marktrs/gitsast/internal/model"
	"github.com/vmihailenco/taskq/v3"
)

var _ IAnalyzeTask = (*Analyzer)(nil)

var cloneLocationPrefix = "temp/"

// Task - register analyze task into task queue
var (
	Task = taskq.RegisterTask(&taskq.TaskOptions{
		Name: "analyzer",
		Handler: func(configPath, id string) error {
			a, err := NewAnalyzer(configPath)
			if err != nil {
				return err
			}

			return a.Analyze(id)
		},
	})
)

// IAnalyzeTask - interface for analyze task
type IAnalyzeTask interface {
	Analyze(reportId string) error
}

type Analyzer struct {
	app    *app.App
	repo   model.IRepositoryRepo
	report model.IReportRepo
	rule   model.IRuleRepo
}

// NewAnalyzeTask - create new analyze task
func NewAnalyzer(configPath string) (IAnalyzeTask, error) {
	rootDir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		return nil, err
	}

	_, app, err := app.Start(context.Background(), "analyzerTask", "", filepath.Join(rootDir, configPath))
	if err != nil {
		return nil, err
	}

	repo := model.NewRepositoryRepo(app)
	report := model.NewReportRepo(app)
	rule := model.NewRuleRepo(app)

	return &Analyzer{app, repo, report, rule}, nil
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
	defer t.removeTempDir(tmpDir)

	log.Infof("scanning files for issues url=%s", repo.RemoteURL)
	issues, err := t.scanFilesForIssues(repo.ID, paths, rules)
	if err != nil {
		return t.handleFailedTask(report, err)
	}

	log.Infof("adding issues to report report_id=%s", report.ID)

	report.Issues = issues
	report.FinishedAt = time.Now()

	if err := t.setReportStatus(report, model.StatusSuccess); err != nil {
		return t.handleFailedTask(report, err)
	}
	log.Infof("analyzed task completed report_id=%s", report.ID)

	return nil
}

// removeTempDir - remove cloned repo directory
func (t *Analyzer) removeTempDir(tmpDir string) error {
	if err := os.RemoveAll(tmpDir); err != nil {
		log.Error(err)
		return err
	}

	return nil
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
