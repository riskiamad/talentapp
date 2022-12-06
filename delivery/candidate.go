package delivery

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"talentapp/model"
	"talentapp/usecase"
	"talentapp/utils"
)

type candidateDelivery struct {
	candidateUsecase usecase.CandidateUsecase
}

func NewCandidateDelivery(candidateUsecase usecase.CandidateUsecase) *candidateDelivery {
	return &candidateDelivery{
		candidateUsecase: candidateUsecase,
	}
}

func (h *candidateDelivery) Router(app *fiber.App) {
	candidate := app.Group("/candidate")
	candidate.Get("/:id", h.GetCandidateByID)
	candidate.Post("", h.PostCandidate)
	candidate.Get("/", h.GetCandidates)
}

func (h *candidateDelivery) GetCandidateByID(ctx *fiber.Ctx) error {
	var (
		id = ctx.Params("id")
	)

	result, err := h.candidateUsecase.GetCandidateByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "something bad happened",
			"error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    result,
	})
}

func (h *candidateDelivery) GetCandidates(ctx *fiber.Ctx) error {
	result, err := h.candidateUsecase.GetCandidates(ctx.Context())
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "something bad happened",
			"error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    result,
	})
}

func (h *candidateDelivery) PostCandidate(ctx *fiber.Ctx) error {
	var (
		payload model.CandidateCreateRequest
		err     error
		ok      bool
		result  = new(model.Candidate)
	)

	if err = ctx.BodyParser(&payload); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request",
			"error":   err.Error(),
		})
	}

	if ok, err = utils.IsRequestValid(payload); !ok {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request",
			"error":   utils.CustomValidator(err),
		})
	}

	result, err = h.candidateUsecase.CreateNewCandidate(ctx.Context(), payload)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "something bad happened",
			"error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"data":    result,
	})
}
