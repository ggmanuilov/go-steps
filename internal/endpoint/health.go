package endpoint

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	startTime = time.Now()
	version   = "1.0.0" // Брать из переменных окружения.
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
}

func Health(c echo.Context) error {
	status := "healthy"

	// Формируем ответ
	response := HealthResponse{
		Status:    status,
		Version:   version,
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).String(),
	}

	if status != "healthy" {
		return c.JSON(http.StatusServiceUnavailable, response)
	}

	return c.JSON(http.StatusOK, response)
}

func Readiness(c echo.Context) error {
	isReady := checkDatabaseConnection()

	if isReady {
		return c.String(http.StatusOK, "Ready")
	}

	return c.String(http.StatusServiceUnavailable, "Not Ready")
}

// Пример функций проверки внешних зависимостей
func checkDatabaseConnection() bool {
	// Пока в базу не ходим. Пока заложено.
	return true
}
