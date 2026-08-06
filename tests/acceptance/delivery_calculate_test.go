package acceptance_test

import (
	"delivery/internal/pkg/app"
	"delivery/internal/requests"
	"delivery/internal/service"
	"delivery/internal/types"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Приемочные тесты расчёта стоимости доставки.
//
// Критерии приёмки (см. README):
//   - Адрес /delivery/calculate рассчитывает доставку с двух складов отправки
//     (Барнаул и Москва) и возвращает самую дешёвую;
//   - ответ содержит склад (Point), стоимость (Cost) и срок (Time) доставки;
//   - если один склад недоступен — используется результат доступного;
//   - если все склады недоступны — возвращается ошибка "delivery unavailable";
//   - невалидные параметры запроса отклоняются с HTTP 400.

// acceptanceFake — фейковый шлюз доставки. Имитирует поведение реального
// шлюза: для каждого запрошенного склада возвращает предзаданный результат,
// в котором Point совпадает с запрошенным складом.
type acceptanceFake struct {
	results map[types.Point]*service.DeliveryResult
}

func (f *acceptanceFake) Calculate(params requests.CalcReq) *service.DeliveryResult {
	if result, ok := f.results[params.PointId]; ok {
		return result
	}
	return &service.DeliveryResult{Error: errors.New("delivery unavailable")}
}

func pointResult(point types.Point, cost float32, time string) *service.DeliveryResult {
	return &service.DeliveryResult{
		CalcResult: types.CalcResult{Point: point, Cost: cost, Time: time},
		Error:      nil,
	}
}

func baseQuery() url.Values {
	q := url.Values{}
	q.Set("DeliveryType", strconv.Itoa(int(types.GatePostal)))
	q.Set("PvzId", "656008")
	q.Set("PointId", strconv.Itoa(int(types.PointBarnaul)))
	q.Set("CountryIso", "643")
	q.Set("Weight", "1.2")
	q.Set("OrderAmount", "100.3")
	return q
}

func requestCalculate(t *testing.T, delivery service.IDelivery, query url.Values) *httptest.ResponseRecorder {
	t.Helper()

	// Собираем приложение через реальную фабрику и реестр маршрутов,
	// чтобы тесты проверяли в том числе проводку /delivery/calculate.
	e := app.New(setDeliveryMiddleware(delivery))
	e.RegisterRoutes()

	req := httptest.NewRequest(http.MethodGet, "/delivery/calculate?"+query.Encode(), nil)
	rec := httptest.NewRecorder()
	e.Echo.ServeHTTP(rec, req)

	return rec
}

// setDeliveryMiddleware внедряет фейковый шлюз доставки в контекст запроса.
func setDeliveryMiddleware(delivery service.IDelivery) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set(string(service.DeliveryServiceKey), delivery)
			return next(ctx)
		}
	}
}

// Возвращается самая дешёвая доставка из доступных складов.
func TestAcceptanceCalculateReturnsCheapestWarehouse(t *testing.T) {
	delivery := &acceptanceFake{results: map[types.Point]*service.DeliveryResult{
		types.PointBarnaul: pointResult(types.PointBarnaul, 90, "3 дня"),
		types.PointMoscow:  pointResult(types.PointMoscow, 110, "2 дня"),
	}}

	rec := requestCalculate(t, delivery, baseQuery())

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"CalcResult":{"Point":1,"Cost":90,"Time":"3 дня"}}`, rec.Body.String())
}

// Самой дешёвой может быть и московская доставка.
func TestAcceptanceCalculateCheapestCanBeMoscow(t *testing.T) {
	delivery := &acceptanceFake{results: map[types.Point]*service.DeliveryResult{
		types.PointBarnaul: pointResult(types.PointBarnaul, 150, "4 дня"),
		types.PointMoscow:  pointResult(types.PointMoscow, 80, "1 день"),
	}}

	rec := requestCalculate(t, delivery, baseQuery())

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"CalcResult":{"Point":2,"Cost":80,"Time":"1 день"}}`, rec.Body.String())
}

// Если один склад недоступен, используется результат доступного.
func TestAcceptanceCalculateFallsBackWhenOneWarehouseUnavailable(t *testing.T) {
	delivery := &acceptanceFake{results: map[types.Point]*service.DeliveryResult{
		types.PointBarnaul: {Error: errors.New("msc unavailable")},
		types.PointMoscow:  pointResult(types.PointMoscow, 120, "2 дня"),
	}}

	rec := requestCalculate(t, delivery, baseQuery())

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"CalcResult":{"Point":2,"Cost":120,"Time":"2 дня"}}`, rec.Body.String())
}

// При равной стоимости доставки из МСК и БРН отправляем с московского склада.
func TestAcceptanceCalculatePrefersMoscowOnEqualCost(t *testing.T) {
	delivery := &acceptanceFake{results: map[types.Point]*service.DeliveryResult{
		types.PointBarnaul: pointResult(types.PointBarnaul, 100, "3 дня"),
		types.PointMoscow:  pointResult(types.PointMoscow, 100, "2 дня"),
	}}

	rec := requestCalculate(t, delivery, baseQuery())

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"CalcResult":{"Point":2,"Cost":100,"Time":"2 дня"}}`, rec.Body.String())
}

// Равная стоимость не отменяет выбор более дешёвого склада.
func TestAcceptanceCalculateEqualCostPreferenceDoesNotOverrideCheaper(t *testing.T) {
	delivery := &acceptanceFake{results: map[types.Point]*service.DeliveryResult{
		types.PointBarnaul: pointResult(types.PointBarnaul, 90, "3 дня"),
		types.PointMoscow:  pointResult(types.PointMoscow, 100, "2 дня"),
	}}

	rec := requestCalculate(t, delivery, baseQuery())

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"CalcResult":{"Point":1,"Cost":90,"Time":"3 дня"}}`, rec.Body.String())
}

// Если все склады недоступны, возвращается ошибка "delivery unavailable".
func TestAcceptanceCalculateAllWarehousesUnavailable(t *testing.T) {
	delivery := &acceptanceFake{results: map[types.Point]*service.DeliveryResult{
		types.PointBarnaul: {Error: errors.New("msc unavailable")},
		types.PointMoscow:  {Error: errors.New("msc unavailable")},
	}}

	rec := requestCalculate(t, delivery, baseQuery())

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"Error":"delivery unavailable"}`, rec.Body.String())
}

// Невалидные параметры отклоняются с HTTP 400.
func TestAcceptanceCalculateRejectsInvalidParameters(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(q url.Values)
	}{
		{"missing Weight", func(q url.Values) { q.Del("Weight") }},
		{"zero Weight", func(q url.Values) { q.Set("Weight", "0") }},
		{"missing OrderAmount", func(q url.Values) { q.Del("OrderAmount") }},
		{"missing PvzId", func(q url.Values) { q.Del("PvzId") }},
		{"missing CountryIso", func(q url.Values) { q.Del("CountryIso") }},
		{"PointId below range", func(q url.Values) { q.Set("PointId", "0") }},
		{"PointId above range", func(q url.Values) { q.Set("PointId", "3") }},
		{"DeliveryType below range", func(q url.Values) { q.Set("DeliveryType", "0") }},
		{"DeliveryType above range", func(q url.Values) { q.Set("DeliveryType", "25") }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			q := baseQuery()
			tc.mutate(q)

			rec := requestCalculate(t, &acceptanceFake{}, q)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), `"Error"`, "ожидался ответ с полем Error")
		})
	}
}

// Эндпоинт /health отвечает статусом healthy.
func TestAcceptanceHealthEndpoint(t *testing.T) {
	e := app.New(nil)
	e.RegisterRoutes()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.Echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"healthy"`)
}
