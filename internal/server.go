package main

import (
	"delivery/internal/pkg/app"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/labstack/gommon/log"
)

func main() {
	handleSigterm()
	app := app.New(nil)
	app.LoadEnv()
	app.Run()
}

// Обработка сигналов завершения работы.
// defer не работает есть получить SIGTERM.
func handleSigterm() {
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		sig := <-c // Ждем код завершения.
		log.Infof("Got OS signal: %s\n", sig.String())

		// Завершаем работу в сервисах.
		// ...

		// Завершаем работу приложения.
		siga, _ := strconv.Atoi(fmt.Sprintf("%d", sig))
		os.Exit(siga)
	}()
}
