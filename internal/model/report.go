package model

import (
	"context"
	"time"

	"github.com/marktrs/gitsast/app"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Report struct {
	bun.BaseModel `bun:"table:reports,alias:report"`

	ID           string       `json:"id" bun:",pk"`
	RepositoryID string       `json:"repository_id" bun:",type:uuid"`
	Status       ReportStatus `json:"status" bun:",default:'initialized'"`
	CreatedAt    time.Time    `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt    time.Time    `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
	EnqueueAt    time.Time    `json:"enqueue_at,omitempty"`
	StartedAt    time.Time    `json:"started_at,omitempty"`
	FinishedAt   time.Time    `json:"finished_at,omitempty"`
	FailedReason string       `json:"failed_reason,omitempty"`

	Issues []*Issue `json:"issues,omitempty" bun:"type:jsonb"`
}

type ReportStatus string

const (
	StatusInitialized ReportStatus = "initialized"
	StatusEnqueued    ReportStatus = "enqueued"
	StatusInProgress  ReportStatus = "in-progress"
	StatusSuccess     ReportStatus = "success"
	StatusFailed      ReportStatus = "failed"
)

// IReport defines methods for read/write reports table.
type IReportRepo interface {
	GetById(ctx context.Context, id string) (*Report, error)
	GetByRepoId(ctx context.Context, id string) (*Report, error)
	Update(ctx context.Context, report *Report) (*Report, error)
	Add(ctx context.Context, report *Report) (*Report, error)
	GetIssues(ctx context.Context, reportID string) ([]*Issue, error)
}

type ReportRepo struct {
	app *app.App
}

func NewReportRepo(app *app.App) IReportRepo {
	return &ReportRepo{app}
}

func (r *ReportRepo) GetById(ctx context.Context, id string) (*Report, error) {
	report := &Report{}
	err := r.app.DB().NewSelect().Model(report).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (r *ReportRepo) GetByRepoId(ctx context.Context, id string) (*Report, error) {
	report := &Report{}
	err := r.app.DB().NewSelect().Model(report).
		Where("repository_id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (r *ReportRepo) Add(ctx context.Context, report *Report) (*Report, error) {
	_, err := r.app.DB().NewInsert().Model(report).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (r *ReportRepo) Update(ctx context.Context, report *Report) (*Report, error) {
	_, err := r.app.DB().NewUpdate().Model(report).WherePK().Exec(ctx)
	if err != nil {
		return nil, err
	}

	return report, nil
}

// GetReportIssues - get all issues for report
func (r *ReportRepo) GetIssues(ctx context.Context, reportID string) ([]*Issue, error) {
	var issues []*Issue
	err := r.app.DB().NewSelect().
		Model((*Report)(nil)).
		ColumnExpr("issues").
		Where("id = ?", reportID).
		Scan(ctx, pgdialect.Array(&issues))
	if err != nil {
		return nil, err
	}

	return issues, nil
}
