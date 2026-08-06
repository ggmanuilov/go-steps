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
	// Приоритет источников окружения:
	//   1. ENV_FILE (задаётся в Dockerfile и make run);
	//   2. .env рядом с исполняемым файлом;
	//   3. .env в рабочем каталоге (покрывает go run и запуск из корня).
	candidates := []string{os.Getenv("ENV_FILE")}
	if execPath, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(execPath), ".env"))
	}
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(cwd, ".env"))
	}

	for _, path := range candidates {
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); err != nil {
			continue
		}
		if err := godotenv.Load(path); err != nil {
			panic("Error loading .env file " + path + ": " + err.Error())
		}
		app.Echo.HideBanner = true
		return
	}

	panic("Error loading .env: neither ENV_FILE nor .env found")
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
