package model

type Candidate struct {
	ID                string `json:"candidate"`
	Name              string `json:"name"`
	Address           string `json:"address"`
	Experience        int    `json:"experience"`
	WillingToRelocate string `json:"willing_to_relocate"`
}

type (
	CandidateCreateRequest struct {
		Name              string `json:"name" validate:"required"`
		Address           string `json:"address" validate:"required"`
		Experience        int    `json:"experience" validate:"required"`
		WillingToRelocate string `json:"willing_to_relocate" validate:"required"`
	}
)
