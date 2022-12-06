package handler

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"talentapp/model"
	"talentapp/usecase"
	"talentapp/utils"
)

type jobHandler struct {
	jobUsecase usecase.JobUsecase
}

func NewJobHandler(jobUsecase usecase.JobUsecase) *jobHandler {
	return &jobHandler{
		jobUsecase: jobUsecase,
	}
}

func (h *jobHandler) Router(app *fiber.App) {
	job := app.Group("/web/job")
	job.Get("", h.Index)
	job.Get("/new", h.New)
	job.Get("/show/:id", h.GetByID)
	job.Post("", h.Create)
}

func (h *jobHandler) Index(ctx *fiber.Ctx) error {
	result, err := h.jobUsecase.GetJobs(ctx.Context())
	if err != nil {
		return ctx.Render("error", nil)
	}

	return ctx.Render(
		"job_index",
		fiber.Map{
			"jobs": result,
		},
	)
}

func (h *jobHandler) New(ctx *fiber.Ctx) error {
	return ctx.Render("job_new", nil)
}

func (h *jobHandler) GetByID(ctx *fiber.Ctx) error {
	var id = ctx.Params("id")

	result, err := h.jobUsecase.GetJobByID(ctx.Context(), id)
	if err != nil {
		return ctx.Render("error", nil)
	}

	return ctx.Render(
		"job_show",
		fiber.Map{
			"job": result,
		},
	)
}

func (h *jobHandler) Create(ctx *fiber.Ctx) error {
	var (
		payload model.JobCreateRequest
		err     error
		ok      bool
	)

	if err = ctx.BodyParser(&payload); err != nil {
		return ctx.Render("job_new", fiber.Map{
			"error": err.Error(),
		})
	}

	if ok, err = utils.IsRequestValid(payload); !ok {
		return ctx.Render("job_new", fiber.Map{
			"error": err.Error(),
		})
	}

	_, err = h.jobUsecase.CreateNewJob(ctx.Context(), payload)
	if err != nil {
		return ctx.Render("job_new", fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Redirect("/web/job", http.StatusFound)
}
