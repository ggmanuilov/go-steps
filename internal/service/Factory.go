package service

import (
	"delivery/internal/requests"
	"delivery/internal/types"
	"fmt"
	"os"
)

type IDelivery interface {
	Calculate(params requests.CalcReq) *DeliveryResult
}

func PointsForCalc() []types.Point {
	return []types.Point{types.PointBarnaul, types.PointMoscow}
}

// DeliveryPool считает с нескольких складов и отдает тот, где дешевле.
func Factory(deliveryType types.DeliveryType, delivery IDelivery) (DeliveryPool, error) {
	// нужно для тестирования с моками реальных вызовов
	if delivery != nil {
		return DeliveryPool{Delivery: delivery}, nil
	}

	switch deliveryType {
	case types.GatePostal:
		postalToken := os.Getenv("POSTAL_TOKEN")
		postalPassword := os.Getenv("POSTAL_PASSWORD")
		if postalToken == "" || postalPassword == "" {
			panic("POSTAL_TOKEN or POSTAL_PASSWORD not set")
		}

		delivery = NewDeliveryPostal(postalToken, postalPassword)
		// case types.CDECK:
		// 	return newDeliveryCdeck()
		// case types.BOXBERRY:
		// 	return newDeliveryBoxberry()
	}

	if delivery == nil {
		return DeliveryPool{Delivery: delivery}, fmt.Errorf("wrong delivery type %d", deliveryType)
	}

	return DeliveryPool{Delivery: delivery}, nil
}
