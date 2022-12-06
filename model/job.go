package model

type Job struct {
	ID             string `json:"id"`
	Position       string `json:"position"`
	Department     string `json:"department"`
	Requester      string `json:"requester"`
	JobDescription string `json:"job_description"`
	Criteria       string `json:"criteria"`
}

type (
	JobCreateRequest struct {
		Position       string `json:"position" validate:"required"`
		Department     string `json:"department" validate:"required"`
		Requester      string `json:"requester" validate:"required"`
		JobDescription string `json:"job_description" validate:"required"`
		Criteria       string `json:"criteria" validate:"required"`
	}
)
