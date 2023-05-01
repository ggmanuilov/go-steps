package delivery

import (
	"errors"
	"sort"
	"sync"
)

type DeliveryPool struct {
	delivery IDelivery
}

type DeliveryResult struct {
	CalcResult CalcResult
	Error      error
}

type ByMinimalCost []DeliveryResult

func (dp ByMinimalCost) Len() int           { return len(dp) }
func (dp ByMinimalCost) Less(i, j int) bool { return dp[i].CalcResult.Cost < dp[j].CalcResult.Cost }
func (dp ByMinimalCost) Swap(i, j int)      { dp[i], dp[j] = dp[j], dp[i] }

func (dp DeliveryPool) Calculate() (CalcResult, error) {
	results := []DeliveryResult{}

	wg := &sync.WaitGroup{}
	for _, pointId := range PointsForCalc() {
		wg.Add(1)
		go func(pointId Point) {
			defer wg.Done()
			results = append(results, dp.delivery.Calculate(pointId))
		}(pointId)
	}
	wg.Wait()

	return dp.FindLessCost(results)
}

func (db DeliveryPool) FindLessCost(results []DeliveryResult) (CalcResult, error) {
	calculated := []DeliveryResult{}
	allIsFail := true
	for _, item := range results {
		if item.Error == nil {
			allIsFail = false
			calculated = append(calculated, item)
		}
	}
	if allIsFail {
		return CalcResult{}, errors.New("Delivery unavailable.")
	}

	sort.Sort(ByMinimalCost(calculated))

	result := calculated[0].CalcResult

	return result, nil
}
