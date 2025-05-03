package integration_test

import (
	"context"
	"delivery/internal/pkg/app"
	"delivery/internal/requests"
	"delivery/internal/service"
	"delivery/internal/types"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createTestServer(middleware echo.MiddlewareFunc) *echo.Echo {
	app := app.New(middleware)

	// Регистрируем маршруты без запуска сервера
	app.RegisterRoutes()

	go func() {
		err := app.Echo.Start(":0")
		if err != nil {
			app.Echo.Logger.Fatal(err)
		}
	}()

	return app.Echo
}

// MockDelivery для интеграционного тестирования
type MockDelivery struct {
	mock.Mock
}

func (m *MockDelivery) Calculate(params requests.CalcReq) *service.DeliveryResult {
	// Предопределенные ответы для конкретных точек доставки
	if params.PointId == types.PointMoscow {
		return &service.DeliveryResult{
			CalcResult: types.CalcResult{
				Cost:  100,
				Time:  "2 дня",
				Point: types.PointBarnaul,
			},
			Error: nil,
		}
	}

	if params.PointId == types.PointBarnaul {
		return &service.DeliveryResult{
			CalcResult: types.CalcResult{
				Cost:  110,
				Time:  "2 дня",
				Point: types.PointMoscow,
			},
			Error: nil,
		}
	}
	args := m.Called(params)
	return args.Get(0).(*service.DeliveryResult)
}

type MockDeliveryUnavailable struct {
	mock.Mock
}

func (m *MockDeliveryUnavailable) Calculate(params requests.CalcReq) *service.DeliveryResult {
	return &service.DeliveryResult{
		CalcResult: types.CalcResult{},
		Error:      fmt.Errorf("delivery unavailable"),
	}
}

// Вспомогательная функция для создания и выполнения запроса
func apiDeliveryRequest(mockDelivery service.IDelivery) *httptest.ResponseRecorder {
	echo := createTestServer(setDeliveryMiddleware(mockDelivery))
	echo.Logger.SetLevel(log.OFF)
	echo.HideBanner = true

	defer echo.Shutdown(context.Background())

	req := httptest.NewRequest(http.MethodGet, "/delivery/calculate", nil)
	q := req.URL.Query()
	q.Add("DeliveryType", strconv.Itoa(int(types.GatePostal)))
	q.Add("PvzId", "656008")
	q.Add("PointId", strconv.Itoa(int(types.PointMoscow)))
	q.Add("CountryIso", "1")
	q.Add("Weight", "1.2")
	q.Add("OrderAmount", "100.3")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	echo.ServeHTTP(rec, req)

	return rec
}

// Вернулась самая дешевая доствка.
func TestDeliveryApiEndpoint(t *testing.T) {
	rec := apiDeliveryRequest(&MockDelivery{})
	assert.Equal(t, http.StatusOK, rec.Code)

	expectedJSON := `{
		"CalcResult": {
			"Point": 1,
			"Cost": 100,
			"Time": "2 дня"
		}
	}`

	assert.JSONEq(t, expectedJSON, rec.Body.String())
}

// Ни один склад не рассчитывает доставку.
func TestDeliveryApiUnavailable(t *testing.T) {
	rec := apiDeliveryRequest(&MockDeliveryUnavailable{})
	assert.Equal(t, http.StatusOK, rec.Code)

	expectedJSON := `{
		"Error": "delivery unavailable"
	}`

	assert.JSONEq(t, expectedJSON, rec.Body.String())
}


func BenchmarkCalculateArray(b *testing.B) {
	deliveryService := service.DeliveryPool{Delivery: &MockDelivery{}}
	for i := 0; i < b.N; i++ {
		deliveryService.Calculate(requests.CalcReq{
			PointId:      types.PointMoscow,
			DeliveryType: types.GatePostal,
			PvzId:        "656008",
			CountryIso:   1,
			Weight:       1.2,
			OrderAmount:  100.3,
		})
	}
}

func setDeliveryMiddleware(mockService service.IDelivery) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// Внедряем мок в контекст
			ctx.Set(string(service.DeliveryServiceKey), mockService)
			return next(ctx)
		}
	}
}
