package endpoint

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	startTime = time.Now()
)

// version возвращает версию сервиса из переменной окружения APP_VERSION
// (задаётся при сборке/деплое), с запасным значением по умолчанию.
// Читается лениво — на каждый запрос, т.к. LoadEnv подгружает .env уже
// после инициализации пакетов.
func version() string {
	if v := os.Getenv("APP_VERSION"); v != "" {
		return v
	}
	return "1.0.0"
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
}

// /health — «процесс жив и HTTP-слой работает» (без зависимостей).
// Отвечает всегда, даже если БД лежит — именно это позволяет отличать «упал сам сервис» от «упала зависимость».
func Health(c echo.Context) error {
	return c.JSON(http.StatusOK, HealthResponse{
		Status:    "healthy",
		Version:   version(),
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).String(),
	})
}

// /readyz — «жив и готов принимать работу» (проверка зависимостей).
func Readiness(c echo.Context) error {
	// Readiness: готовность работать вместе с внешними зависимостями.
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
