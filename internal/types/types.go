package types

type Point uint8

const (
	PointBarnaul Point = 1
	PointMoscow  Point = 2
)

type CalcParams struct {
	PointId     uint8 // склад отправки посылки
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

type PostalCalcResp struct {
	Cost float32
	Time string
}

type DeliveryType int8

const (
	GatePostal   DeliveryType = 1
	GateBoxberry DeliveryType = 25
)
