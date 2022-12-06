package model

import (
	"context"
	"talentapp/driver/db/mysql"
	"time"
)

type Recruitment struct {
	ID             string    `json:"id"`
	JobID          string    `json:"job_id"`
	Job            *Job      `json:"job,omitempty"`
	Status         string    `json:"status"`
	Deadline       time.Time `json:"deadline"`
	DeadlineString string    `json:"-"`
}

type (
	RecruitmentCreateRequest struct {
		JobID    string `json:"job_id" validate:"required"`
		Deadline string `json:"deadline" validate:"required"`
	}

	RecruitmentUpdateStatusRequest struct {
		Status string `json:"status" validate:"required"`
	}
)

func (r *Recruitment) LoadJob(ctx context.Context) error {
	var result Job
	SQL := "SELECT id, position, department, requester, job_description, criteria FROM job WHERE id = ?"
	row := mysql.DB.QueryRowContext(ctx, SQL, r.JobID)

	err := row.Scan(&result.ID, &result.Position, &result.Department, &result.Requester, &result.JobDescription, &result.Criteria)
	if err != nil {
		return err
	}

	r.Job = &result

	return nil
}
