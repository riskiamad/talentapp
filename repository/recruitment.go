package repository

import (
	"context"
	"database/sql"
	"strings"
	"talentapp/model"
)

type RecruitmentRepository interface {
	GetRecruitments(ctx context.Context) (*[]model.Recruitment, error)
	GetRecruitmentByID(ctx context.Context, id string) (*model.Recruitment, error)
	CreateRecruitment(ctx context.Context, model *model.Recruitment) error
	UpdateRecruitmentStatus(ctx context.Context, model *model.Recruitment) error
}

type recruitmentRepository struct {
	DB *sql.DB
}

func NewRecruitmentRepository(db *sql.DB) RecruitmentRepository {
	return &recruitmentRepository{
		db,
	}
}

func (r *recruitmentRepository) GetRecruitmentByID(ctx context.Context, id string) (*model.Recruitment, error) {
	var result model.Recruitment
	SQL := "SELECT id, job_id, status, deadline FROM recruitment WHERE id = ?"
	row := r.DB.QueryRowContext(ctx, SQL, id)

	err := row.Scan(&result.ID, &result.JobID, &result.Status, &result.Deadline)
	if err != nil {
		return nil, err
	}

	result.DeadlineString = strings.Split(result.Deadline.String(), " ")[0]

	return &result, nil
}

func (r *recruitmentRepository) CreateRecruitment(ctx context.Context, model *model.Recruitment) error {
	SQL := "insert into recruitment(id, job_id, status, deadline) values (?, ?, ?, ?)"
	if _, err := r.DB.ExecContext(ctx, SQL, model.ID, model.JobID, model.Status, model.Deadline); err != nil {
		return err
	}

	return nil
}

func (r *recruitmentRepository) GetRecruitments(ctx context.Context) (*[]model.Recruitment, error) {
	var result []model.Recruitment
	SQL := "SELECT id, job_id, status, deadline FROM recruitment LIMIT ? OFFSET ?"
	rows, err := r.DB.QueryContext(ctx, SQL, 10, 0) // TODO: make dynamic pagination
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		recruitment := model.Recruitment{}
		err = rows.Scan(&recruitment.ID, &recruitment.JobID, &recruitment.Status, &recruitment.Deadline)
		if err != nil {
			return nil, err
		}

		err = recruitment.LoadJob(ctx)
		if err != nil {
			return nil, err
		}

		recruitment.DeadlineString = strings.Split(recruitment.Deadline.String(), " ")[0]

		result = append(result, recruitment)
	}

	return &result, nil
}

func (r *recruitmentRepository) UpdateRecruitmentStatus(ctx context.Context, model *model.Recruitment) error {
	SQL := "update recruitment set status = ? where id = ?"
	if _, err := r.DB.ExecContext(ctx, SQL, model.Status, model.ID); err != nil {
		return err
	}

	return nil
}
