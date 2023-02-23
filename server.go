package main

import (
	"fmt"
	"net/http"
	"os"
	"shop-cart/db"
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

	database, err := db.NewDatabase(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASSWORD"))
	if err != nil {
		log.Sugar().Errorf("Failed to connect to redis: %s", err.Error())
	}
	defer database.Client.Close()

	e.GET("/", func(c echo.Context) error {
		c.Logger().Debug('/')
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":" + os.Getenv("APP_PORT")))
}
