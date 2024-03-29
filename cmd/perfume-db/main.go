package main

import (
	"context"
	"fmt"
	"github.com/ej-agas/perfume-db/internal"
	"github.com/ej-agas/perfume-db/nanoid"
	"github.com/ej-agas/perfume-db/postgresql"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

type config struct {
	port          int
	environment   string
	encryptionKey string
}

type application struct {
	config    config
	logger    *slog.Logger
	validator *validator.Validate
	services  *postgresql.Services
	factory   *internal.Factory
}

var Version string

func main() {
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))

	if err != nil {
		log.Fatal(fmt.Errorf("invalid application port: %s", err))
	}

	cfg := config{
		port:          port,
		environment:   os.Getenv("APP_ENV"),
		encryptionKey: os.Getenv("APP_ENCRYPTION_KEY"),
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal(fmt.Errorf("invalid database port: %s", err))
	}

	connString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		dbPort,
		os.Getenv("DB_NAME"),
	)

	connConfig, err := pgx.ParseConfig(connString)
	if err != nil {
		panic(err)
	}

	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	idAlphabet := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	idLength := 16
	idGenerator := nanoid.NewNanoIdGenerator(idAlphabet, idLength)

	app := &application{
		config:    cfg,
		logger:    slog.New(slog.NewTextHandler(os.Stderr, nil)),
		validator: validator.New(validator.WithRequiredStructEnabled()),
		services:  postgresql.NewServices(conn),
		factory:   &internal.Factory{IdGenerator: idGenerator},
	}

	app.logger.Info("APP RUNNING IN", "PORT", os.Getenv("APP_PORT"))

	app.logger.Error(http.ListenAndServe(":"+os.Getenv("APP_PORT"), app.routes()).Error())
}
