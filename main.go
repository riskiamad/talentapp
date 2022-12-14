package main

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
	"github.com/golang-migrate/migrate/v4"
	mysqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/mattes/migrate/source/file"
	"log"
	"os"
	"talentapp/delivery"
	"talentapp/driver/db/mysql"
	"talentapp/handler"
	"talentapp/repository"
	"talentapp/usecase"
)

var db *sql.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	db = mysql.ConnectDB()

	dbInstance, err := mysqlMigrate.WithInstance(db, &mysqlMigrate.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		dbInstance,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err = m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}
}

func main() {
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
