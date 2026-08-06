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

	res, err := sharedClient.Do(req)
	if err != nil {
		return &gateCalcResp, fmt.Errorf("postal gate: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return &gateCalcResp, fmt.Errorf("postal gate: unexpected status %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &gateCalcResp, fmt.Errorf("postal gate: read body: %w", err)
	}

	err = json.Unmarshal(data, &gateCalcResp)
	if err != nil {
		return &gateCalcResp, fmt.Errorf("postal gate: unmarshal: %w", err)
	}

	return &gateCalcResp, nil
}

// sharedClient — общий HTTP-клиент без собственного таймаута: дедлайн задаёт
// контекст запроса, чтобы был один источник правды о таймауте.
var sharedClient = &http.Client{}

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
			Cost:  rand.Float32(), // todo это заглушка
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
