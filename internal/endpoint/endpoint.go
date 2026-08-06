package endpoint

import (
	"delivery/internal/requests"
	"delivery/internal/service"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	deliveryPool *service.DeliveryPool
}

func New(deliveryPool service.DeliveryPool) *Endpoint {
	return &Endpoint{
		deliveryPool: &deliveryPool,
	}
}

func validate(c echo.Context) (*requests.CalcReq, error) {
	// добавляем валидацию стркутуры
	requests.CalcRegister(c.Echo())
	calcReq := new(requests.CalcReq)
	if err := c.Bind(calcReq); err != nil {
		return calcReq, err
	}

	if err := c.Validate(calcReq); err != nil {
		return calcReq, err
	}

	return calcReq, nil
}

func Calculate(c echo.Context) error {

	calcReq, err := validate(c)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"Error": err.Error(),
		})
	}

	deliveryPool, err := GetDeliveryService(c, calcReq)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"Error": err.Error(),
		})
	}

	// считаем доставку
	calcResult := deliveryPool.Calculate(*calcReq)
	if calcResult.Error != nil {
		return c.JSON(http.StatusOK, map[string]any{
			"Error": calcResult.Error.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"CalcResult": calcResult.CalcResult,
	})
}

func GetDeliveryService(c echo.Context, calcReq *requests.CalcReq) (service.DeliveryPool, error) {
	var deliveryPool service.DeliveryPool
	var err error
	// в тестах будем мокать вызовы к сторонним сервисам
	delivery := c.Get(string(service.DeliveryServiceKey))
	if delivery != nil {
		deliveryService, ok := delivery.(service.IDelivery)
		if !ok {
			return deliveryPool, fmt.Errorf("delivery service not found in context")

		}
		deliveryPool, err = service.Factory(calcReq.DeliveryType, deliveryService)
	} else {
		deliveryPool, err = service.Factory(calcReq.DeliveryType, nil)
	}

	return deliveryPool, err
}
