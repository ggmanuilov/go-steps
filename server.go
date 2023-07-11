package main

import (
	"fmt"
	"net/http"
	"os"
	"shop-cart/delivery"
	"shop-cart/requests"
	"shop-cart/utils"

	"github.com/brpaz/echozap"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	log, err := utils.InitializeLogger()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	e := echo.New()
	e.Use(echozap.ZapLogger(log))

	err = godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}

	// Learn connect redis to project.
	// database, err := db.NewDatabase(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASSWORD"))
	// if err != nil {
	// 	log.Sugar().Errorf("Failed to connect to redis: %s", err.Error())
	// }
	// defer database.Client.Close()

	e.GET("/", func(c echo.Context) error {
		c.Logger().Debug('/')
		return c.String(http.StatusOK, "Hello, World!")
	})

	requests.CalcRegister(e)

	e.GET("/delivery/calculate", func(c echo.Context) error {
		// валидация
		calcReq := new(requests.CalcReq)
		if err = c.Bind(calcReq); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		if err = c.Validate(calcReq); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		// расчет стоимости
		c.Logger().SetPrefix("/delivery/calculate")

		params := delivery.CalcParams{
			CountryIso:  calcReq.CountryIso,
			GateId:      calcReq.GateId,
			Weight:      calcReq.Weight,
			OrderAmount: calcReq.OrderAmount,
		}

		deliveryPool, err := delivery.Factory(calcReq.Type, params)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		calcResult, err := deliveryPool.Calculate()
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, calcResult)
	})

	e.Logger.Fatal(e.Start(":" + os.Getenv("APP_PORT")))
}
