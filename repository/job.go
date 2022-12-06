package repository

import (
	"context"
	"database/sql"
	"talentapp/model"
)

type JobRepository interface {
	GetJobByID(ctx context.Context, id string) (*model.Job, error)
	PostJob(ctx context.Context, model *model.Job) error
	GetJobList(ctx context.Context) (*[]model.Job, error)
}

type jobRepository struct {
	DB *sql.DB
}

func NewJobRepository(db *sql.DB) JobRepository {
	return &jobRepository{DB: db}
}

func (r *jobRepository) GetJobByID(ctx context.Context, id string) (*model.Job, error) {
	var result model.Job
	SQL := "SELECT id, position, department, requester, job_description, criteria FROM job WHERE id = ?"
	row := r.DB.QueryRowContext(ctx, SQL, id)

	err := row.Scan(&result.ID, &result.Position, &result.Department, &result.Requester, &result.JobDescription, &result.Criteria)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *jobRepository) PostJob(ctx context.Context, model *model.Job) error {
	SQL := "insert into job(id, position, department, requester, job_description, criteria) values (?, ?, ?, ?, ?, ?)"
	if _, err := r.DB.ExecContext(ctx, SQL, model.ID, model.Position, model.Department, model.Requester, model.JobDescription, model.Criteria); err != nil {
		return err
	}

	return nil
}

func (r *jobRepository) GetJobList(ctx context.Context) (*[]model.Job, error) {
	var result []model.Job
	SQL := "SELECT id, position, department, requester, job_description, criteria FROM job LIMIT ? OFFSET ?"
	rows, err := r.DB.QueryContext(ctx, SQL, 10, 0) // TODO: make dynamic pagination
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		job := model.Job{}
		err = rows.Scan(&job.ID, &job.Position, &job.Department, &job.Requester, &job.JobDescription, &job.Criteria)
		if err != nil {
			return nil, err
		}
		result = append(result, job)
	}

	return &result, nil
}
