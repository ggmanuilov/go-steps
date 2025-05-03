package service

import (
	"context"
	"delivery/internal/requests"
	"delivery/internal/types"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type PostalGate struct {
	accessToken    string
	loginPassword  string
	timeoutSeconds int64
}

func (d PostalGate) callApi(params requests.CalcReq) (*types.PostalCalcResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(d.timeoutSeconds)*time.Second)
	defer cancel()

	var gateCalcResp types.PostalCalcResp

	log.Printf("Gate: %s, Params: %+v\n", params.PvzId, params)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://otpravka.pochta.ru/1.0/tariff", nil)
	if err != nil {
		return &gateCalcResp, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json;charset=UTF-8")
	req.Header.Set("Authorization", "AccessToken "+d.accessToken)
	req.Header.Set("X-User-Authorization", "Basic "+d.loginPassword)

	c := http.Client{Timeout: time.Duration(1) * time.Second}

	res, err := c.Do(req)
	if err != nil {
		return &gateCalcResp, err
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &gateCalcResp)
	if err != nil {
		return &gateCalcResp, err
	}

	return &gateCalcResp, nil
}

func (d PostalGate) Calculate(params requests.CalcReq) *DeliveryResult {

	gateCalcResp, err := d.callApi(params)
	if err != nil {
		resp := DeliveryResult{
			CalcResult: types.CalcResult{},
			Error:      errors.New("msc unavailable"),
		}
		return &resp
	}

	return &DeliveryResult{
		CalcResult: types.CalcResult{
			Point: types.Point(params.PointId),
			Time:  fmt.Sprintf("%v days", gateCalcResp.Time),
			Cost:  rand.Float32(),
		},
		Error: nil,
	}
}

func NewDeliveryPostal(accessToken string, loginPassword string) IDelivery {
	return PostalGate{
		accessToken:    accessToken,
		loginPassword:  loginPassword,
		timeoutSeconds: 20,
	}
}
