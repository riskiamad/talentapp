package repository

import (
	"context"
	"database/sql"
	"talentapp/model"
)

type CandidateRepository interface {
	GetCandidateByID(ctx context.Context, id string) (*model.Candidate, error)
	PostCandidate(ctx context.Context, model *model.Candidate) error
	GetCandidateList(ctx context.Context) (*[]model.Candidate, error)
}

type candidateRepository struct {
	DB *sql.DB
}

func NewCandidateRepository(db *sql.DB) CandidateRepository {
	return &candidateRepository{DB: db}
}

func (r *candidateRepository) GetCandidateByID(ctx context.Context, id string) (*model.Candidate, error) {
	var result model.Candidate
	SQL := "SELECT id, name, address, experience, willing_to_relocate FROM candidate WHERE id = ?"
	row := r.DB.QueryRowContext(ctx, SQL, id)

	err := row.Scan(&result.ID, &result.Name, &result.Address, &result.Experience, &result.WillingToRelocate)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *candidateRepository) PostCandidate(ctx context.Context, model *model.Candidate) error {
	SQL := "insert into candidate(id, name, address, experience, willing_to_relocate) values (?, ?, ?, ?, ?)"
	if _, err := r.DB.ExecContext(ctx, SQL, model.ID, model.Name, model.Address, model.Experience, model.WillingToRelocate); err != nil {
		return err
	}

	return nil
}

func (r *candidateRepository) GetCandidateList(ctx context.Context) (*[]model.Candidate, error) {
	var result []model.Candidate
	SQL := "SELECT id, name, address, experience, willing_to_relocate FROM candidate LIMIT ? OFFSET ?"
	rows, err := r.DB.QueryContext(ctx, SQL, 10, 0) // TODO: make dynamic pagination
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		candidate := model.Candidate{}
		err = rows.Scan(&candidate.ID, &candidate.Name, &candidate.Address, &candidate.Experience, &candidate.WillingToRelocate)
		if err != nil {
			return nil, err
		}
		result = append(result, candidate)
	}

	return &result, nil
}
