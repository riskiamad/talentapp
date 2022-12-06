package repository

import (
	"context"
	"database/sql"
	"talentapp/model"
)

type CandidateScoreRepository interface {
	GetCandidateScoreByID(ctx context.Context, id, candidateID string) (*model.CandidateScore, error)
	PostCandidateScore(ctx context.Context, model *model.CandidateScore) error
	GetCandidateScoreListByRecruitmentID(ctx context.Context, recruitmentID string) (*[]model.CandidateScore, error)
}

type candidateScoreRepository struct {
	DB *sql.DB
}

func NewCandidateScoreRepository(db *sql.DB) CandidateScoreRepository {
	return &candidateScoreRepository{DB: db}
}

func (r *candidateScoreRepository) GetCandidateScoreByID(ctx context.Context, id, candidateID string) (*model.CandidateScore, error) {
	var result model.CandidateScore
	SQL := "SELECT id, candidate_id, recruitment_id, willing_to_relocate_score, attitude_score, skill_score, experience_score, overall_score FROM candidate_score WHERE recruitment_id = ? AND candidate_id = ?"
	rows, err := r.DB.QueryContext(ctx, SQL, id, candidateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rows.Scan(&result.ID, &result.CandidateID, &result.RecruitmentID, &result.WillingToRelocateScore, &result.AttitudeScore, &result.SkillScore, &result.ExperienceScore, &result.OverallScore)

	return &result, nil
}

func (r *candidateScoreRepository) PostCandidateScore(ctx context.Context, model *model.CandidateScore) error {
	SQL := "insert into candidate_score(id, candidate_id, recruitment_id, willing_to_relocate_score, attitude_score, skill_score, experience_score, overall_score) values (?, ?, ?, ?, ?, ?, ?, ?)"
	if _, err := r.DB.ExecContext(ctx, SQL, model.ID, model.CandidateID, model.RecruitmentID, model.WillingToRelocateScore, model.AttitudeScore, model.SkillScore, model.ExperienceScore, model.OverallScore); err != nil {
		return err
	}

	return nil
}

func (r *candidateScoreRepository) GetCandidateScoreListByRecruitmentID(ctx context.Context, recruitmentID string) (*[]model.CandidateScore, error) {
	var result []model.CandidateScore
	SQL := "SELECT id, candidate_id, recruitment_id, willing_to_relocate_score, attitude_score, skill_score, experience_score, overall_score FROM candidate_score WHERE recruitment_id = ? ORDER BY ? DESC LIMIT ? OFFSET ?"
	rows, err := r.DB.QueryContext(ctx, SQL, recruitmentID, "overall_score", 10, 0) // TODO: make dynamic pagination
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		candidateScore := model.CandidateScore{}
		err = rows.Scan(&candidateScore.ID, &candidateScore.CandidateID, &candidateScore.RecruitmentID, &candidateScore.WillingToRelocateScore, &candidateScore.AttitudeScore, &candidateScore.SkillScore, &candidateScore.ExperienceScore, &candidateScore.OverallScore)
		if err != nil {
			return nil, err
		}

		err = candidateScore.LoadCandidate(ctx)
		if err != nil {
			return nil, err
		}

		err = candidateScore.LoadRecruitment(ctx)
		if err != nil {
			return nil, err
		}

		err = candidateScore.Recruitment.LoadJob(ctx)
		if err != nil {
			return nil, err
		}

		result = append(result, candidateScore)
	}

	return &result, nil
}
