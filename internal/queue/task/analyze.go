package task

import (
	"context"
	"database/sql"
	"os"

	"github.com/labstack/gommon/log"
	"github.com/marktrs/gitsast/internal/config"
	"github.com/marktrs/gitsast/internal/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/vmihailenco/taskq/v3"
)

var _ IAnalyzeTask = (*Analyzer)(nil)

var (
	AnalyzeTask = taskq.RegisterTask(&taskq.TaskOptions{
		Name: "analyzer",
		Handler: func(id string) error {
			NewAnalyzeTask().Start(id)
			return nil
		},
	})
)

type IAnalyzeTask interface {
	Start(reportId string) error
}

type Analyzer struct {
	repo   model.IRepository
	report model.IReport
}

func NewAnalyzeTask() IAnalyzeTask {
	// TODO: don't reload config here
	// load app config
	cfg, err := config.Load("./config/local.yaml")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// connect db
	db := bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DB.Dsn))), pgdialect.New())
	repo := model.NewRepositoryRepo(db)
	report := model.NewReportRepo(db)

	return &Analyzer{repo, report}
}

func (t *Analyzer) Start(reportId string) error {
	ctx := context.Background()

	// set status
	report, err := t.report.GetById(ctx, reportId)
	if err != nil {
		return err
	}

	log.Infof("starting analyzed task report_id=%s", report.ID)
	if err := t.setReportStatus(report, model.StatusInProgress); err != nil {
		return err
	}

	// 1. look up for rules

	// 2. look up for latest code

	// 3. scan for issue

	// 4. add issue to report (db)

	// 5. update report status (db)

	if err := t.setReportStatus(report, model.StatusSuccess); err != nil {
		return err
	}

	// 6. end task
	log.Infof("analyzed task completed report_id=%s", report.ID)

	return nil
}

func (t *Analyzer) setReportStatus(report *model.Report, status model.ReportStatus) error {
	ctx := context.Background()

	report.Status = status

	_, err := t.report.Update(ctx, report)
	if err != nil {
		return err
	}

	return nil
}
