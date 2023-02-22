package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/marktrs/gitsast/internal/model"
	"github.com/marktrs/gitsast/internal/queue"
	"github.com/marktrs/gitsast/internal/queue/task/analyzer"
)

var _ IService = (*service)(nil)

var (
	ErrReportInProgress = errors.New(
		`the report for this repository already initialized, only completed/failed report can retry`)
)

// IService defines methods for business logic of repository domain
// such as validate request body, enqueue analyzing task, CRUD repository or report
type IService interface {
	GetById(ctx context.Context, id string) (*model.Repository, error)
	List(ctx context.Context, f *model.RepositoryFilter) ([]*model.Repository, error)
	Add(ctx context.Context, req *AddRepositoryRequest) (*model.Repository, error)
	Update(ctx context.Context, id string, req *UpdateRepositoryRequest) error
	Remove(ctx context.Context, id string) error
	CreateReport(ctx context.Context, repoId string) (*model.Report, error)
	GetReportByRepoId(ctx context.Context, repoId string) (*model.Report, error)
}

type service struct {
	repo      model.IRepositoryRepo
	report    model.IReportRepo
	queue     queue.Handler
	validator *validator.Validate
}

type AddRepositoryRequest struct {
	Name      string `json:"name" validate:"required,max=120"`
	RemoteURL string `json:"remote_url" validate:"required,max=120"`
}

type UpdateRepositoryRequest struct {
	Name      string `json:"name"`
	RemoteURL string `json:"remote_url"`
}

func NewService(v *validator.Validate, q queue.Handler, rs model.IRepositoryRepo, rp model.IReportRepo) IService {
	return &service{
		repo:      rs,
		report:    rp,
		validator: v,
		queue:     q,
	}
}

// GetById implements IService.GetById interface.
func (s *service) GetById(ctx context.Context, id string) (*model.Repository, error) {
	return s.repo.GetById(ctx, id)
}

// List implements IService.List interface.
func (s *service) List(ctx context.Context, f *model.RepositoryFilter) ([]*model.Repository, error) {
	return s.repo.List(ctx, f)
}

// Add implements IService.Add interface.
func (s *service) Add(ctx context.Context, req *AddRepositoryRequest) (*model.Repository, error) {
	// validate request body
	if err := s.validator.Struct(req); err != nil {
		log.Errorf("Request validation failed on add repository handler : %s", err.Error())
		return nil, err
	}

	repo := &model.Repository{
		ID:        uuid.New().String(),
		Name:      req.Name,
		RemoteURL: req.RemoteURL,
	}

	return s.repo.Add(ctx, repo)
}

// Update implements IService.Update interface.
func (s *service) Update(ctx context.Context, id string, req *UpdateRepositoryRequest,
) error {
	// validate request body
	if err := s.validator.Struct(req); err != nil {
		log.Errorf("Request validation failed on update repository handler : %s", err.Error())
		return err
	}

	repo := map[string]interface{}{
		"name":       req.Name,
		"remote_url": req.RemoteURL,
		"updated_at": time.Now(),
	}

	return s.repo.Update(ctx, id, repo)
}

// Remove implements IService.Remove interface.
func (s *service) Remove(ctx context.Context, id string) error {
	return s.repo.Remove(ctx, id)
}

// CreateReport implements IService.CreateReport interface.
func (s *service) CreateReport(ctx context.Context, repoId string) (*model.Report, error) {
	repo, err := s.repo.GetById(ctx, repoId)
	if err != nil {
		return nil, err
	}

	report, err := s.report.GetByRepoId(ctx, repo.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	// if report exist
	if report != nil {
		if report.Status != model.StatusSuccess && report.Status != model.StatusFailed {
			return nil, ErrReportInProgress
		}
	}

	// generate a new report
	report = &model.Report{
		ID:           uuid.New().String(),
		Status:       model.StatusInitialized,
		RepositoryID: repoId,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Issues:       []*model.Issue{},
	}

	report, err = s.report.Add(ctx, report)
	if err != nil {
		return nil, err
	}
	log.Infof("created a new report with id: %s", report.ID)

	s.queue.AddTask(
		analyzer.Task.WithArgs(ctx, report.ID))
	if err != nil {
		return nil, err
	}
	log.Infof("enqueued a new report with id: %s", report.ID)

	report.Status = model.StatusEnqueued
	report.EnqueueAt = time.Now()

	report, err = s.report.Update(ctx, report)
	if err != nil {
		return nil, err
	}

	return report, nil
}

// GetReport - Implements IService.GetReport interface.
func (s *service) GetReportByRepoId(ctx context.Context, id string) (*model.Report, error) {
	return s.report.GetByRepoId(ctx, id)
}
