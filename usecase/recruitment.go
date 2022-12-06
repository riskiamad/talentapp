package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"math"
	"talentapp/model"
	"talentapp/repository"
	"talentapp/utils"
	"time"
)

type RecruitmentUsecase interface {
	CreateNewRecruitment(ctx context.Context, payload model.RecruitmentCreateRequest) (*model.Recruitment, error)
	GetRecruitmentByID(ctx context.Context, id string) (*model.Recruitment, error)
	GetRecruitments(ctx context.Context) (*[]model.Recruitment, error)
	UpdateRecruitmentStatus(ctx context.Context, id string, payload model.RecruitmentUpdateStatusRequest) (*model.Recruitment, error)
	GetRecruitmentScores(ctx context.Context, id string) (*[]model.CandidateScore, error)
	CreateNewCandidateScore(ctx context.Context, recruitmentID string, payload model.CandidateScoreCreateRequest) (*model.CandidateScore, error)
	GetCandidateScoreByID(ctx context.Context, id, candidateID string) (*model.CandidateScore, error)
}

type recruitmentUsecase struct {
	recruitmentRepository    repository.RecruitmentRepository
	candidateScoreRepository repository.CandidateScoreRepository
	candidateRepository      repository.CandidateRepository
	jobRepository            repository.JobRepository
}

func NewRecruitmentUsecase(
	recruitmentRepository repository.RecruitmentRepository,
	candidateScoreRepository repository.CandidateScoreRepository,
	candidateRepository repository.CandidateRepository,
	jobRepository repository.JobRepository,
) RecruitmentUsecase {
	return &recruitmentUsecase{
		recruitmentRepository:    recruitmentRepository,
		candidateScoreRepository: candidateScoreRepository,
		candidateRepository:      candidateRepository,
		jobRepository:            jobRepository,
	}
}

func (u *recruitmentUsecase) GetRecruitmentByID(ctx context.Context, id string) (*model.Recruitment, error) {
	result, err := u.recruitmentRepository.GetRecruitmentByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("recruitment with id %s not found", id)
	}

	if err != nil {
		return nil, err
	}

	if err = result.LoadJob(ctx); err != nil {
		return nil, err
	}

	return result, nil
}

func (u *recruitmentUsecase) GetRecruitments(ctx context.Context) (*[]model.Recruitment, error) {
	result, err := u.recruitmentRepository.GetRecruitments(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *recruitmentUsecase) CreateNewRecruitment(ctx context.Context, payload model.RecruitmentCreateRequest) (*model.Recruitment, error) {
	var (
		result = new(model.Recruitment)
		newID  = uuid.NewString()
	)

	deadline, err := time.Parse("2006-01-02", payload.Deadline)
	if err != nil {
		return nil, err
	}

	result = &model.Recruitment{
		ID:       newID,
		JobID:    payload.JobID, // TODO: validate job id
		Status:   "open",
		Deadline: deadline.Add(time.Hour * 23).Add(time.Minute * 59).Add(time.Second * 59),
	}

	err = u.recruitmentRepository.CreateRecruitment(ctx, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *recruitmentUsecase) UpdateRecruitmentStatus(ctx context.Context, id string, payload model.RecruitmentUpdateStatusRequest) (*model.Recruitment, error) {
	var (
		result = new(model.Recruitment)
	)

	result, err := u.recruitmentRepository.GetRecruitmentByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("recruitment with id %s not found", id)
	}

	if err != nil {
		return nil, err
	}

	result.Status = payload.Status

	err = u.recruitmentRepository.UpdateRecruitmentStatus(ctx, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *recruitmentUsecase) GetRecruitmentScores(ctx context.Context, id string) (*[]model.CandidateScore, error) {
	result, err := u.candidateScoreRepository.GetCandidateScoreListByRecruitmentID(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *recruitmentUsecase) CreateNewCandidateScore(ctx context.Context, recruitmentID string, payload model.CandidateScoreCreateRequest) (*model.CandidateScore, error) {
	var (
		result                                                             = new(model.CandidateScore)
		newID                                                              = uuid.NewString()
		willingToRelocateScore, experienceScore, attitudeScore, skillScore float64
	)

	switch payload.WillingToRelocate {
	case "yes":
		willingToRelocateScore = 100
	case "no":
		willingToRelocateScore = 50
	default:
		willingToRelocateScore = 50
	}

	switch exp := payload.Experience; {
	case exp == 1:
		experienceScore = 50
	case exp == 2:
		experienceScore = 70
	case exp == 3:
		experienceScore = 80
	case exp == 4:
		experienceScore = 90
	case exp >= 5:
		experienceScore = 100
	default:
		experienceScore = 50
	}

	attitudeScore = utils.GetFloat64Value(payload.AttitudeScore)
	skillScore = utils.GetFloat64Value(payload.SkillScore)

	attitude := attitudeScore * 0.2
	willingToRelocate := willingToRelocateScore * 0.2
	skill := skillScore * 0.3
	experience := experienceScore * 0.3
	overallScore := math.Round((attitude+willingToRelocate+skill+experience)*100) / 100

	result = &model.CandidateScore{
		ID:                     newID,
		CandidateID:            payload.CandidateID, // TODO: validate candidateID
		RecruitmentID:          recruitmentID,       // TODO: validate recruitmentID
		WillingToRelocateScore: willingToRelocateScore,
		AttitudeScore:          attitudeScore,
		SkillScore:             skillScore,
		ExperienceScore:        experienceScore,
		OverallScore:           overallScore,
	}

	err := u.candidateScoreRepository.PostCandidateScore(ctx, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *recruitmentUsecase) GetCandidateScoreByID(ctx context.Context, id, candidateID string) (*model.CandidateScore, error) {
	result, err := u.candidateScoreRepository.GetCandidateScoreByID(ctx, id, candidateID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("candidate score with id %s not found", id)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}
