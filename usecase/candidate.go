package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"talentapp/model"
	"talentapp/repository"
)

type CandidateUsecase interface {
	CreateNewCandidate(ctx context.Context, payload model.CandidateCreateRequest) (*model.Candidate, error)
	GetCandidateByID(ctx context.Context, id string) (*model.Candidate, error)
	GetCandidates(ctx context.Context) (*[]model.Candidate, error)
}

type candidateUsecase struct {
	candidateRepository repository.CandidateRepository
}

func NewCandidateUsecase(candidateRepository repository.CandidateRepository) CandidateUsecase {
	return &candidateUsecase{
		candidateRepository: candidateRepository,
	}
}

func (u *candidateUsecase) GetCandidateByID(ctx context.Context, id string) (*model.Candidate, error) {
	result, err := u.candidateRepository.GetCandidateByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("candidate with id %s not found", id)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *candidateUsecase) GetCandidates(ctx context.Context) (*[]model.Candidate, error) {
	result, err := u.candidateRepository.GetCandidateList(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *candidateUsecase) CreateNewCandidate(ctx context.Context, payload model.CandidateCreateRequest) (*model.Candidate, error) {
	var (
		result = new(model.Candidate)
		newID  = uuid.NewString()
	)

	result = &model.Candidate{
		ID:                newID,
		Name:              payload.Name,
		Address:           payload.Address,
		Experience:        payload.Experience,
		WillingToRelocate: payload.WillingToRelocate,
	}

	err := u.candidateRepository.PostCandidate(ctx, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
