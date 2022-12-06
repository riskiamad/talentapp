package handler

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"talentapp/model"
	"talentapp/usecase"
	"talentapp/utils"
)

type candidateHandler struct {
	candidateUsecase usecase.CandidateUsecase
}

func NewCandidateHandler(candidateUsecase usecase.CandidateUsecase) *candidateHandler {
	return &candidateHandler{
		candidateUsecase: candidateUsecase,
	}
}

func (h *candidateHandler) Router(app *fiber.App) {
	candidate := app.Group("/web/candidate")
	candidate.Get("", h.Index)
	candidate.Get("/new", h.New)
	candidate.Get("/show/:id", h.GetByID)
	candidate.Post("", h.Create)
}

func (h *candidateHandler) Index(ctx *fiber.Ctx) error {
	result, err := h.candidateUsecase.GetCandidates(ctx.Context())
	if err != nil {
		return ctx.Render("error", nil)
	}

	return ctx.Render(
		"candidate_index",
		fiber.Map{
			"candidates": result,
		},
	)
}

func (h *candidateHandler) New(ctx *fiber.Ctx) error {
	return ctx.Render("candidate_new", nil)
}

func (h *candidateHandler) GetByID(ctx *fiber.Ctx) error {
	var id = ctx.Params("id")

	result, err := h.candidateUsecase.GetCandidateByID(ctx.Context(), id)
	if err != nil {
		return ctx.Render("error", nil)
	}

	return ctx.Render(
		"candidate_show",
		fiber.Map{
			"candidate": result,
		},
	)
}

func (h *candidateHandler) Create(ctx *fiber.Ctx) error {
	var (
		payload model.CandidateCreateRequest
		err     error
		ok      bool
	)

	if err = ctx.BodyParser(&payload); err != nil {
		return ctx.Render("candidate_new", fiber.Map{
			"error": err,
		})
	}

	if ok, err = utils.IsRequestValid(payload); !ok {
		return ctx.Render("candidate_new", fiber.Map{
			"error": err,
		})
	}

	_, err = h.candidateUsecase.CreateNewCandidate(ctx.Context(), payload)
	if err != nil {
		return ctx.Render("candidate_new", fiber.Map{
			"error": err,
		})
	}

	return ctx.Redirect("/web/candidate", http.StatusFound)
}
