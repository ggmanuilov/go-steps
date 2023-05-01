package delivery

import (
	"errors"
	"testing"
)

// Gate unavailable.
func TestAllIsFail(t *testing.T) {
	results := []DeliveryResult{}
	results = append(results, DeliveryResult{
		CalcResult: CalcResult{},
		Error:      errors.New("Api unavailable"),
	})
	results = append(results, DeliveryResult{
		CalcResult: CalcResult{},
		Error:      errors.New("Api unavailable"),
	})

	dp := DeliveryPool{
		delivery: newDeliveryPostal(CalcParams{}),
	}

	_, gotError := dp.FindLessCost(results)
	wantError := errors.New("Delivery unavailable.")
	if gotError.Error() != wantError.Error() {
		t.Errorf("wrong error: got %v, want %v", gotError, wantError)
	}
}
