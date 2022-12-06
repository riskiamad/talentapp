package delivery

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"talentapp/model"
	"talentapp/usecase"
	"talentapp/utils"
)

type recruitmentDelivery struct {
	recruitmentUsecase usecase.RecruitmentUsecase
}

func NewRecruitmentDelivery(recruitmentUsecase usecase.RecruitmentUsecase) *recruitmentDelivery {
	return &recruitmentDelivery{
		recruitmentUsecase: recruitmentUsecase,
	}
}

func (h *recruitmentDelivery) Router(app *fiber.App) {
	recruitment := app.Group("/recruitment")
	recruitment.Get("/:id", h.GetRecruitmentByID)
	recruitment.Post("", h.PostRecruitment)
	recruitment.Get("/", h.GetRecruitments)
	recruitment.Put("/:id", h.PutRecruitmentStatus)
	recruitment.Get("/:id/score", h.GetRecruitmentScore)
	recruitment.Post("/:id/score", h.PostCandidateScore)
	recruitment.Get("/:id/candidate/:candidate_id/score", h.GetCandidateScoreByID)
}

func (h *recruitmentDelivery) GetRecruitmentByID(ctx *fiber.Ctx) error {
	var (
		id = ctx.Params("id")
	)

	result, err := h.recruitmentUsecase.GetRecruitmentByID(ctx.Context(), id)
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

func (h *recruitmentDelivery) GetRecruitments(ctx *fiber.Ctx) error {
	result, err := h.recruitmentUsecase.GetRecruitments(ctx.Context())
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

func (h *recruitmentDelivery) PostRecruitment(ctx *fiber.Ctx) error {
	var (
		payload model.RecruitmentCreateRequest
		err     error
		ok      bool
		result  = new(model.Recruitment)
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

	result, err = h.recruitmentUsecase.CreateNewRecruitment(ctx.Context(), payload)
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

func (h *recruitmentDelivery) PutRecruitmentStatus(ctx *fiber.Ctx) error {
	var (
		id      = ctx.Params("id")
		payload model.RecruitmentUpdateStatusRequest
		err     error
		ok      bool
		result  = new(model.Recruitment)
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

	result, err = h.recruitmentUsecase.UpdateRecruitmentStatus(ctx.Context(), id, payload)
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

func (h *recruitmentDelivery) GetRecruitmentScore(ctx *fiber.Ctx) error {
	var id = ctx.Params("id")

	result, err := h.recruitmentUsecase.GetRecruitmentScores(ctx.Context(), id)
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

func (h *recruitmentDelivery) PostCandidateScore(ctx *fiber.Ctx) error {
	var (
		payload model.CandidateScoreCreateRequest
		id      = ctx.Params("id")
		err     error
		ok      bool
		result  = new(model.CandidateScore)
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

	result, err = h.recruitmentUsecase.CreateNewCandidateScore(ctx.Context(), id, payload)
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

func (h *recruitmentDelivery) GetCandidateScoreByID(ctx *fiber.Ctx) error {
	var (
		id          = ctx.Params("id")
		candidateID = ctx.Params("candidate_id")
	)

	result, err := h.recruitmentUsecase.GetCandidateScoreByID(ctx.Context(), id, candidateID)
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
