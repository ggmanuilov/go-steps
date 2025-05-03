package delivery

import (
	"context"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type postal struct {
	params CalcParams
	client *http.Client
}

func (d postal) Calculate(pointId Point, client *http.Client) (*DeliveryResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*100))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://postal.ru/any-params", nil)
	if err != nil {
		return nil, err
	}

	c := http.Client{Timeout: time.Duration(1) * time.Second}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	out, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(out))
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

	resp := DeliveryResult{
		CalcResult: result,
		Error:      nil,
	}

	return &resp, nil
}

func newDeliveryPostal(calcParams CalcParams, client *http.Client) IDelivery {
	return postal{
		params: calcParams,
		client: client,
	}
}
