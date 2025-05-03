package app

import (
	"delivery/internal/endpoint"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type App struct {
	Echo *echo.Echo
}

func New(middleware echo.MiddlewareFunc) *App {
	app := &App{
		Echo: echo.New(),
	}

	if middleware != nil {
		app.Echo.Use(middleware)
	}

	return app
}

func (app *App) LoadEnv() {
	execPath, err := os.Executable()
	if err != nil {
		panic("Error getting executable path: " + err.Error())
	}

	execDir := filepath.Dir(execPath)
	err = godotenv.Load(filepath.Join(execDir, ".env"))
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}
	app.Echo.HideBanner = true
}

func (app *App) RegisterRoutes() {
	app.Echo.GET("/health", endpoint.Health)
	app.Echo.GET("/readyz", endpoint.Readiness)

	app.Echo.GET("/delivery/calculate", endpoint.Calculate)
}

func (app *App) Run() {

	app.RegisterRoutes()

	err := app.Echo.Start(":" + os.Getenv("APP_PORT"))
	if err != nil {
		app.Echo.Logger.Fatal(err)
	}
}
