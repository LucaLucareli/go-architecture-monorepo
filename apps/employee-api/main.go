package main

import (
	"employee-api/modules"
	"fmt"
	"os"
	"shared"
	"shared/pkg/helpers"
	"shared/pkg/interceptors"
	"shared/pkg/logger"
	"shared/pkg/validation"
	exceptionfactory "shared/pkg/validation/exception-factory"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	godotenv.Load()

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

	logger.Init("EMPLOYEE-API", logger.ColorBlue, "DEV")

	modules.NewAppModule(e, appState)

	port := helpers.GetEnv("PORT", 3000)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
