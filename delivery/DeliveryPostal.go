package delivery

import (
	"math/rand"
	"time"
)

type postal struct {
	params CalcParams
}

func (d postal) getCost() (float32, error) {
	return 0, nil
}

func (d postal) getTime() (string, error) {
	return "0 days", nil
}

func (d postal) Calculate(pointId Point) (deliveryResult DeliveryResult) {

	// Simulate when one delivery point unavailable.
	/*
		if pointId == PointMoscow {
			return DeliveryResult{
				CalcResult: CalcResult{},
				Error:      errors.New("Msc unavailable"),
			}
		}*/

	// Simulate
	time.Sleep(100*time.Millisecond - time.Duration(pointId))

	result := CalcResult{
		Point: pointId,
		Time:  "2 days",
		Cost:  rand.Float32(),
	}

	return DeliveryResult{
		CalcResult: result,
		Error:      nil,
	}
}

func newDeliveryPostal(calcParams CalcParams) IDelivery {
	return &postal{
		params: calcParams,
	}
}
