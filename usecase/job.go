package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"talentapp/model"
	"talentapp/repository"
)

type JobUsecase interface {
	CreateNewJob(ctx context.Context, payload model.JobCreateRequest) (*model.Job, error)
	GetJobByID(ctx context.Context, id string) (*model.Job, error)
	GetJobs(ctx context.Context) (*[]model.Job, error)
}

type jobUsecase struct {
	jobRepository repository.JobRepository
}

func NewJobUsecase(jobRepository repository.JobRepository) JobUsecase {
	return &jobUsecase{
		jobRepository,
	}
}

func (u *jobUsecase) GetJobByID(ctx context.Context, id string) (*model.Job, error) {
	result, err := u.jobRepository.GetJobByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("job with id %s not found", id)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *jobUsecase) GetJobs(ctx context.Context) (*[]model.Job, error) {
	result, err := u.jobRepository.GetJobList(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *jobUsecase) CreateNewJob(ctx context.Context, payload model.JobCreateRequest) (*model.Job, error) {
	var (
		result = new(model.Job)
		newID  = uuid.NewString()
	)

	result = &model.Job{
		ID:             newID,
		Position:       payload.Position,
		Department:     payload.Department,
		Requester:      payload.Requester,
		JobDescription: payload.JobDescription,
		Criteria:       payload.Criteria,
	}

	err := u.jobRepository.PostJob(ctx, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
