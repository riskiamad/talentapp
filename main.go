package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	"log"
	"os"
	"talentapp/delivery"
	"talentapp/driver/db/mysql"
	"talentapp/handler"
	"talentapp/repository"
	"talentapp/usecase"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {
	db := mysql.ConnectDB()

	engine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(cors.New())
	app.Static("/images", "./images")
	app.Static("/css", "./assets/css")
	app.Static("/js", "./assets/js")
	app.Static("/webfonts", "./assets/webfonts")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to talentapp")
	})

	// dependencies
	// repository
	jobRepository := repository.NewJobRepository(db)
	candidateScoreRepository := repository.NewCandidateScoreRepository(db)
	recruitmentRepository := repository.NewRecruitmentRepository(db)
	candidateRepository := repository.NewCandidateRepository(db)

	// usecase
	jobUsecase := usecase.NewJobUsecase(jobRepository)
	recruitmentUsecase := usecase.NewRecruitmentUsecase(
		recruitmentRepository,
		candidateScoreRepository,
		candidateRepository,
		jobRepository,
	)
	candidateUsecase := usecase.NewCandidateUsecase(candidateRepository)

	// delivery
	jobDelivery := delivery.NewJobDelivery(jobUsecase)
	recruitmentDelivery := delivery.NewRecruitmentDelivery(recruitmentUsecase)
	candidateDelivery := delivery.NewCandidateDelivery(candidateUsecase)

	// handler
	recruitmentHandler := handler.NewRecruitmentHandler(recruitmentUsecase, jobUsecase, candidateUsecase)
	jobHandler := handler.NewJobHandler(jobUsecase)
	candidateHandler := handler.NewCandidateHandler(candidateUsecase)

	// router
	jobDelivery.Router(app)
	recruitmentDelivery.Router(app)
	candidateDelivery.Router(app)
	recruitmentHandler.Router(app)
	jobHandler.Router(app)
	candidateHandler.Router(app)

	err := app.Listen(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}
}
