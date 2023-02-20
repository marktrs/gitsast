package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"
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
type IReport interface {
	GetById(ctx context.Context, id string) (*Report, error)
	GetByRepoId(ctx context.Context, id string) (*Report, error)
	Update(ctx context.Context, report *Report) (*Report, error)
	Add(ctx context.Context, report *Report) (*Report, error)
}

type report struct {
	db *bun.DB
}

func NewReportRepo(db *bun.DB) IReport {
	return &report{db}
}

func (r *report) GetById(ctx context.Context, id string) (*Report, error) {
	report := &Report{}
	err := r.db.NewSelect().Model(report).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (r *report) GetByRepoId(ctx context.Context, id string) (*Report, error) {
	report := &Report{}
	err := r.db.NewSelect().Model(report).
		Where("repository_id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (r *report) Add(ctx context.Context, report *Report) (*Report, error) {
	_, err := r.db.NewInsert().Model(report).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (r *report) Update(ctx context.Context, report *Report) (*Report, error) {
	_, err := r.db.NewUpdate().Model(report).WherePK().Exec(ctx)
	if err != nil {
		return nil, err
	}

	return report, nil
}
