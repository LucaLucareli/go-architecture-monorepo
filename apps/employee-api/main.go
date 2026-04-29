package main

import (
	"employee-api/modules"
	"fmt"
	"os"
	"time"

	"shared"
	"shared/pkg/helpers"
	"shared/pkg/interceptors"
	"shared/pkg/logger"
	"shared/pkg/middlewares"
	"shared/pkg/validation"
	exceptionfactory "shared/pkg/validation/exception-factory"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system environment variables")
	}

	dbPostgresURL := os.Getenv("DATABASE_URL")
	dbRedisURL := os.Getenv("REDIS_URL")

	appState := shared.NewAppState(dbPostgresURL, dbRedisURL)

	e := echo.New()

	e.Validator = validation.NewValidator()
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if handled := exceptionfactory.CustomExceptionFactory(err); handled != nil {
			c.JSON(handled.Code, handled.Message)
			return
		}
		e.DefaultHTTPErrorHandler(err, c)
	}

	e.Use(interceptors.RequestLogger)
	e.Use(interceptors.TransformInterceptor)
	e.Use(middlewares.TimeoutMiddleware(30 * time.Second))
	e.Use(middlewares.AsyncAuditMiddleware())
	e.Use(middlewares.RateLimitMiddleware(20, 100))

	logger.Init("EMPLOYEE-API", logger.ColorBlue, "DEV")

	appModule := modules.NewAppModule()
	appModule.RegisterAllRoutes(e, appState)

	port := helpers.GetEnv("PORT", 3000)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
