package delivery

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"talentapp/model"
	"talentapp/usecase"
	"talentapp/utils"
)

type jobDelivery struct {
	jobUsecase usecase.JobUsecase
}

func NewJobDelivery(jobUsecase usecase.JobUsecase) *jobDelivery {
	return &jobDelivery{
		jobUsecase: jobUsecase,
	}
}

func (h *jobDelivery) Router(app *fiber.App) {
	job := app.Group("/job")
	job.Get("/:id", h.GetJobByID)
	job.Post("", h.PostJob)
	job.Get("", h.GetJobs)
}

func (h *jobDelivery) GetJobByID(ctx *fiber.Ctx) error {
	var (
		id = ctx.Params("id")
	)

	result, err := h.jobUsecase.GetJobByID(ctx.Context(), id)
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

func (h *jobDelivery) GetJobs(ctx *fiber.Ctx) error {
	result, err := h.jobUsecase.GetJobs(ctx.Context())
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

func (h *jobDelivery) PostJob(ctx *fiber.Ctx) error {
	var (
		payload model.JobCreateRequest
		err     error
		ok      bool
		result  = new(model.Job)
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

	result, err = h.jobUsecase.CreateNewJob(ctx.Context(), payload)
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
