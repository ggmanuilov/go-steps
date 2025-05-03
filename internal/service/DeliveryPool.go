package service

import (
	"delivery/internal/requests"
	"delivery/internal/types"
	"errors"
	"sort"
	"sync"
)

type DeliveryPool struct {
	Delivery IDelivery
}

type DeliveryResult struct {
	CalcResult types.CalcResult `json:",omitempty"`
	Error      error            `json:",omitempty"`
}

type ByMinimalCost []DeliveryResult

func (dp ByMinimalCost) Len() int           { return len(dp) }
func (dp ByMinimalCost) Less(i, j int) bool { return dp[i].CalcResult.Cost < dp[j].CalcResult.Cost }
func (dp ByMinimalCost) Swap(i, j int)      { dp[i], dp[j] = dp[j], dp[i] }

func (dp DeliveryPool) CalculateOne(params requests.CalcReq) *DeliveryResult {
	return dp.Delivery.Calculate(params)
}

func (dp DeliveryPool) Calculate(params requests.CalcReq) *DeliveryResult {
	var results [2]DeliveryResult

	wg := &sync.WaitGroup{}
	var mu sync.Mutex
	var resultsCount int

	for _, pointId := range PointsForCalc() {
		wg.Add(1)
		go func(pointId types.Point) {
			defer wg.Done()

			params.PointId = pointId
			result := dp.Delivery.Calculate(params)

			mu.Lock()
			if resultsCount < 2 {
				results[resultsCount] = *result
				resultsCount++
			}
			mu.Unlock()
		}(pointId)
	}
	wg.Wait()

	return dp.findLessCost(results[:2])
}

func (db DeliveryPool) findLessCost(results []DeliveryResult) *DeliveryResult {
	calculated := make([]DeliveryResult, 0, len(results))
	allIsFail := true
	for _, item := range results {
		if item.Error == nil {
			allIsFail = false
			calculated = append(calculated, item)
		}
	}

	if allIsFail {
		return &DeliveryResult{
			CalcResult: types.CalcResult{},
			Error:      errors.New("delivery unavailable"),
		}
	}

	sort.Sort(ByMinimalCost(calculated))

	result := calculated[0].CalcResult

	return &DeliveryResult{
		CalcResult: result,
		Error:      nil,
	}
}
