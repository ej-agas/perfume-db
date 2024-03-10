package main

import (
	"fmt"
	"github.com/ej-agas/perfume-db/handlers"
	"github.com/ej-agas/perfume-db/postgresql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

type config struct {
	port        int
	environment string
}

type application struct {
	config       config
	logger       *slog.Logger
	houseHandler *handlers.HouseHandler
}

var Version string

func main() {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))

	if err != nil {
		log.Fatal(fmt.Errorf("invalid DB port: %s", err))
	}

	cfg := config{
		port:        port,
		environment: os.Getenv("APP_ENV"),
	}

	app := &application{
		config: cfg,
		logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}

	app.houseHandler = &handlers.HouseHandler{Service: postgresql.HouseService{}}

	app.logger.Info("APP RUNNING IN", "PORT", os.Getenv("APP_PORT"))

	app.logger.Error(http.ListenAndServe(":"+os.Getenv("APP_PORT"), app.routes()).Error())
}
