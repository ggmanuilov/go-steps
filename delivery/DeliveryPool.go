package delivery

import (
	"errors"
)

type DeliveryPool struct {
	delivery IDelivery
}

type DeliveryResult struct {
	CalcResult CalcResult
	Error      error
}

func (dp DeliveryPool) Calculate() (CalcResult, error) {
	// Есть два склада отправки.
	// Нужно рассчитать каждый и вернуть тот, где дешевле стоимость доставки.
	// Если стоимость доставки одинакова, вернуть склад PointMoscow.
	// Если расчет недоступен со всех складов вернуть ошибку.

	chBrnRes := make(chan DeliveryResult)
	chMscRes := make(chan DeliveryResult)

	defer close(chBrnRes)
	defer close(chMscRes)

	// Расчет доставки здесь.
	go dp.delivery.Calculate(PointBarnaul, chBrnRes)
	go dp.delivery.Calculate(PointMoscow, chMscRes)

	calcResultMsc := <-chMscRes
	calcResultBrn := <-chBrnRes

	return dp.getLessCost(calcResultMsc, calcResultBrn)
}

func (db DeliveryPool) getLessCost(calcResultMsc DeliveryResult, calcResultBrn DeliveryResult) (CalcResult, error) {
	errMsk := calcResultMsc.Error
	errBrn := calcResultMsc.Error

	if errMsk != nil && errBrn != nil {
		var result CalcResult
		return result, errors.New("Delivery unavailable.")
	}

	if errMsk != nil && errBrn == nil {
		return calcResultBrn.CalcResult, nil
	}

	if errMsk == nil && errBrn != nil {
		return calcResultMsc.CalcResult, nil
	}

	if calcResultMsc.CalcResult.Cost <= calcResultBrn.CalcResult.Cost {
		return calcResultMsc.CalcResult, nil
	}

	return calcResultBrn.CalcResult, nil
}
