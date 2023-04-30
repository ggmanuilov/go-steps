package delivery

import (
	"fmt"
)

const (
	TypePostal       int8 = 1
	TypePostalOnline int8 = 14
	TypeSdec         int8 = 11
	TypeBoxberry     int8 = 25
	TypeEuroset      int8 = 2
)

type Point uint8

const (
	PointBarnaul Point = 1
	PointMoscow  Point = 2
)

func PointsForCalc() []uint8 {
	return []uint8{uint8(PointBarnaul), uint8(PointMoscow)}
}

type CalcParams struct {
	CountryIso  uint16
	GateId      string
	Weight      float32
	OrderAmount float32
}

type CalcResult struct {
	Point Point   // склад отправки посылки
	Cost  float32 // стоимость доставки
	Time  string  // время доставки
}

type IDelivery interface {
	getCost() (float32, error)
	getTime() (string, error)
	Calculate(pointId Point, calcDeliveryResult chan DeliveryResult)
}

// DeliveryPool считает с нескольких складов и отдает тот, где дешевле.
func Factory(deliveryType int8, calcParams CalcParams) (DeliveryPool, error) {
	var delivery IDelivery
	switch deliveryType {
	case TypePostal:
		delivery = newDeliveryPostal(calcParams)
		// case POSTAL_ONLINE:
		// 	return newDeliveryPostalOnline()
		// case CDECK:
		// 	return newDeliveryCdeck()
		// case BOXBERRY:
		// 	return newDeliveryBoxberry()
		// case EUROSET:
		// 	return newEuroset()
	}

	if delivery == nil {
		return DeliveryPool{delivery: delivery}, fmt.Errorf("Wrong delivery type")
	}

	return DeliveryPool{delivery: delivery}, nil
}
