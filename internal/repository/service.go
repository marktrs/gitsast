package repository

import (
	"context"
	"database/sql"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/marktrs/gitsast/app"
	"github.com/marktrs/gitsast/internal/model"
	"github.com/marktrs/gitsast/internal/queue"
	"github.com/marktrs/gitsast/internal/queue/task/analyzer"
	"github.com/rs/zerolog/log"
)

// IService variable that does static check to make sure that 'service' struct implements 'IService' interface.
var _ IService = (*service)(nil)

// IService defines methods for business logic of repository domain
// such as validate request body, enqueue analyzing task, CRUD repository or report
type IService interface {
	GetById(ctx context.Context, id string) (*model.Repository, error)
	List(ctx context.Context, f *model.RepositoryFilter) ([]*model.Repository, error)
	Add(ctx context.Context, req *AddRepositoryRequest) (*model.Repository, error)
	Update(ctx context.Context, id string, req *UpdateRepositoryRequest) error
	Remove(ctx context.Context, id string) error
	CreateReport(ctx context.Context, repoId string) (*model.Report, error)
	GetReportByRepoId(ctx context.Context, repoId string) (*GetReportResponse, error)
}

type service struct {
	app *app.App

	repo      model.IRepositoryRepo
	report    model.IReportRepo
	queue     queue.Handler
	validator *validator.Validate
}

type AddRepositoryRequest struct {
	Name      string `json:"name" validate:"required,max=120"`
	RemoteURL string `json:"remote_url" validate:"required,max=120,is-git-url"`
}

func (r *AddRepositoryRequest) Validate(validator *validator.Validate) error {
	return validator.Struct(r)
}

type UpdateRepositoryRequest struct {
	Name      string `json:"name" validate:"max=120"`
	RemoteURL string `json:"remote_url" validate:"max=120,is-git-url"`
}

func (r *UpdateRepositoryRequest) Validate(validator *validator.Validate) error {
	return validator.Struct(r)
}

func NewService(app *app.App, rs model.IRepositoryRepo, rp model.IReportRepo) IService {
	app.Validator().RegisterValidation("is-git-url", ValidateGitRemoteURL)
	return &service{
		app:       app,
		repo:      rs,
		report:    rp,
		queue:     app.Queue(),
		validator: app.Validator(),
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
	if err := req.Validate(s.validator); err != nil {
		log.Err(err).Msg("request validation failed on add repository handler")
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
func (s *service) Update(
	ctx context.Context,
	id string,
	req *UpdateRepositoryRequest,
) error {
	// validate request body
	if err := req.Validate(s.validator); err != nil {
		log.Err(err).Msg("request validation failed on update repository handler")
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
			return nil, model.ErrReportInProgress
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

	// enqueue a new analyzing task to main queue
	err = s.queue.AddTask(analyzer.Task.WithArgs(ctx, report.ID))
	if err != nil {
		return nil, err
	}

	report.Status = model.StatusEnqueued
	report.EnqueueAt = time.Now()

	report, err = s.report.Update(ctx, report)
	if err != nil {
		return nil, err
	}

	return report, nil
}

type GetReportResponse struct {
	model.Report
	Findings []*model.Finding `json:"findings"`
}

// GetReport - Implements IService.GetReport interface.
func (s *service) GetReportByRepoId(ctx context.Context, id string) (*GetReportResponse, error) {
	report, err := s.report.GetByRepoId(ctx, id)
	if err != nil {
		return nil, err
	}

	var findings []*model.Finding
	for _, issue := range report.Issues {
		var finding model.Finding
		finding.Type = "sast"
		finding.RuleID = issue.RuleID
		finding.Location.Path = issue.Location.Path
		finding.Location.Position.Begin.Line = int(issue.Location.Line)
		finding.Metadata.Description = issue.Description
		finding.Metadata.Severity = issue.Severity
		findings = append(findings, &finding)
	}

	var response GetReportResponse
	response.Report = *report
	response.Issues = nil
	response.Findings = findings

	return &response, nil
}

func ValidateGitRemoteURL(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	if url == "" {
		return false
	}

	re := regexp.MustCompile(`(?P<Protocol>git@|http(s)?:\/\/)(.+@)*(?P<Provider>[\w\d\.-]+)(:[\d]+){0,1}(\/scm)?\/*(?P<Name>.*)`)
	matches := re.FindStringSubmatch(url)

	return len(matches) != 0
}
