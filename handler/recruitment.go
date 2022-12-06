package handler

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"talentapp/model"
	"talentapp/usecase"
	"talentapp/utils"
	"time"
)

type recruitmentHandler struct {
	recruitmentUsecase usecase.RecruitmentUsecase
	jobUsecase         usecase.JobUsecase
	candidateUsecase   usecase.CandidateUsecase
}

func NewRecruitmentHandler(
	recruitmentUsecase usecase.RecruitmentUsecase,
	jobUsecase usecase.JobUsecase,
	candidateUsecase usecase.CandidateUsecase,
) *recruitmentHandler {
	return &recruitmentHandler{
		recruitmentUsecase: recruitmentUsecase,
		jobUsecase:         jobUsecase,
		candidateUsecase:   candidateUsecase,
	}
}

func (h *recruitmentHandler) Router(app *fiber.App) {
	recruitment := app.Group("/web/recruitment")
	recruitment.Get("", h.Index)
	recruitment.Get("/new", h.New)
	recruitment.Get("/show/:id", h.GetByID)
	recruitment.Post("", h.Create)
	recruitment.Get("/show/:id/score", h.ScoreList)
	recruitment.Get("/:id/score/new", h.NewScore)
	recruitment.Post("/edit/:id/status", h.UpdateStatus)
	recruitment.Post("/:id/score", h.CreateScore)
}

func (h *recruitmentHandler) Index(ctx *fiber.Ctx) error {
	result, err := h.recruitmentUsecase.GetRecruitments(ctx.Context())
	if err != nil {
		return ctx.Render("error", nil)
	}

	return ctx.Render(
		"recruitment_index",
		fiber.Map{
			"recruitments": result,
		},
	)
}

func (h *recruitmentHandler) New(ctx *fiber.Ctx) error {
	now := time.Now().Format("2006-01-02")
	jobs, err := h.jobUsecase.GetJobs(ctx.Context())
	if err != nil {
		return ctx.Render("error", nil)
	}

	return ctx.Render(
		"recruitment_new",
		fiber.Map{
			"jobs": jobs,
			"now":  now,
		})
}

func (h *recruitmentHandler) GetByID(ctx *fiber.Ctx) error {
	var id = ctx.Params("id")

	result, err := h.recruitmentUsecase.GetRecruitmentByID(ctx.Context(), id)
	if err != nil {
		return ctx.Render("error", nil)
	}

	return ctx.Render(
		"recruitment_show",
		fiber.Map{
			"recruitment": result,
		},
	)
}

func (h *recruitmentHandler) Create(ctx *fiber.Ctx) error {
	var (
		payload model.RecruitmentCreateRequest
		err     error
		ok      bool
	)

	if err = ctx.BodyParser(&payload); err != nil {
		return ctx.Render("recruitment_new", fiber.Map{
			"error": err.Error(),
		})
	}

	if ok, err = utils.IsRequestValid(payload); !ok {
		return ctx.Render("recruitment_new", fiber.Map{
			"error": err.Error(),
		})
	}

	_, err = h.recruitmentUsecase.CreateNewRecruitment(ctx.Context(), payload)
	if err != nil {
		return ctx.Render("recruitment_new", fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Redirect("/web/recruitment", http.StatusFound)
}

func (h *recruitmentHandler) ScoreList(ctx *fiber.Ctx) error {
	var id = ctx.Params("id")

	result, err := h.recruitmentUsecase.GetRecruitmentScores(ctx.Context(), id)
	if err != nil {
		return ctx.Render("error", nil)
	}

	return ctx.Render(
		"score_index",
		fiber.Map{
			"scores":        result,
			"recruitmentID": id,
		},
	)
}

func (h *recruitmentHandler) NewScore(ctx *fiber.Ctx) error {
	var id = ctx.Params("id")

	candidates, err := h.candidateUsecase.GetCandidates(ctx.Context())
	if err != nil {
		return ctx.Render("error", nil)
	}

	return ctx.Render(
		"score_new",
		fiber.Map{
			"candidates":    candidates,
			"recruitmentID": id,
		})
}

func (h *recruitmentHandler) UpdateStatus(ctx *fiber.Ctx) error {
	var (
		id      = ctx.Params("id")
		payload model.RecruitmentUpdateStatusRequest
		ok      bool
		err     error
	)

	result, err := h.recruitmentUsecase.GetRecruitments(ctx.Context())
	if err != nil {
		return ctx.Render("error", nil)
	}

	if err = ctx.BodyParser(&payload); err != nil {
		return ctx.Render("recruitment_index", fiber.Map{
			"error":        err.Error(),
			"recruitments": result,
		})
	}

	if ok, err = utils.IsRequestValid(payload); !ok {
		return ctx.Render("recruitment_index", fiber.Map{
			"error":        err.Error(),
			"recruitments": result,
		})
	}

	if payload.Status == "open" {
		payload.Status = "close"
	} else if payload.Status == "close" {
		payload.Status = "open"
	}

	h.recruitmentUsecase.UpdateRecruitmentStatus(ctx.Context(), id, payload)
	if err != nil {
		return ctx.Render("recruitment_index", fiber.Map{
			"error":        err.Error(),
			"recruitments": result,
		})
	}

	return ctx.Redirect("/web/recruitment", http.StatusFound)
}

func (h *recruitmentHandler) CreateScore(ctx *fiber.Ctx) error {
	var (
		payload model.CandidateScoreCreateRequest
		err     error
		ok      bool
		id      = ctx.Params("id")
	)

	candidates, err := h.candidateUsecase.GetCandidates(ctx.Context())
	if err != nil {
		return ctx.Render("error", nil)
	}

	if err = ctx.BodyParser(&payload); err != nil {
		return ctx.Render("recruitment_new", fiber.Map{
			"error":         err.Error(),
			"candidates":    candidates,
			"recruitmentID": id,
		})
	}

	if ok, err = utils.IsRequestValid(payload); !ok {
		return ctx.Render("score_new", fiber.Map{
			"error":         err.Error(),
			"candidates":    candidates,
			"recruitmentID": id,
		})
	}

	_, err = h.recruitmentUsecase.CreateNewCandidateScore(ctx.Context(), id, payload)
	if err != nil {
		return ctx.Render("score_new", fiber.Map{
			"error":         err.Error(),
			"candidates":    candidates,
			"recruitmentID": id,
		})
	}

	return ctx.Redirect("/web/recruitment/show/"+id+"/score", http.StatusFound)
}
