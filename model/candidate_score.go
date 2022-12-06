package model

import (
	"context"
	"talentapp/driver/db/mysql"
)

type CandidateScore struct {
	ID                     string       `json:"id"`
	CandidateID            string       `json:"candidate_id"`
	Candidate              *Candidate   `json:"candidate,omitempty"`
	RecruitmentID          string       `json:"recruitment_id"`
	Recruitment            *Recruitment `json:"recruitment,omitempty"`
	WillingToRelocateScore float64      `json:"willing_to_relocate_score"`
	AttitudeScore          float64      `json:"attitude_score"`
	SkillScore             float64      `json:"skill_score"`
	ExperienceScore        float64      `json:"experience_score"`
	OverallScore           float64      `json:"overall_score"`
}

type (
	CandidateScoreCreateRequest struct {
		CandidateID       string `json:"candidate_id" validate:"required"`
		WillingToRelocate string `json:"willing_to_relocate_score" validate:"required,eq=yes|eq=no"`
		AttitudeScore     string `json:"attitude_score" validate:"required"`
		SkillScore        string `json:"skill_score" validate:"required"`
		Experience        int    `json:"experience" validate:"required"`
	}
)

func (cs *CandidateScore) LoadCandidate(ctx context.Context) error {
	var result Candidate
	SQL := "SELECT id, name, address, experience, willing_to_relocate FROM candidate WHERE id = ?"
	row := mysql.DB.QueryRowContext(ctx, SQL, cs.CandidateID)

	err := row.Scan(&result.ID, &result.Name, &result.Address, &result.Experience, &result.WillingToRelocate)
	if err != nil {
		return err
	}

	cs.Candidate = &result

	return nil
}

func (cs *CandidateScore) LoadRecruitment(ctx context.Context) error {
	var result Recruitment
	SQL := "SELECT id, job_id, status, deadline FROM recruitment WHERE id = ?"
	row := mysql.DB.QueryRowContext(ctx, SQL, cs.RecruitmentID)

	err := row.Scan(&result.ID, &result.JobID, &result.Status, &result.Deadline)
	if err != nil {
		return err
	}

	cs.Recruitment = &result

	return nil
}
