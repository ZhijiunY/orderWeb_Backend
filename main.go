package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/ZhijiunY/orderWeb_Backend/models"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const VERSION = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	jwt struct {
		secret string
	}
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type application struct {
	config config
	logger *log.Logger
	models models.Models
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 5001, "Server Port To Listen On")
	flag.StringVar(&cfg.env, "env", "development", "Application Environment (development|production)")
	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://fang@localhost/afternoon-tea?sslmode=disable", "Postgres Connection String")
	flag.StringVar(&cfg.jwt.secret, "jwt-secret", "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160", "secret")
	flag.Parse()

	// cfg.jwt.secret = os.Getenv("GO_AFTERNOON_JWT")

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	server := &http.Server{
		Addr:         fmt.Sprint(":", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("STARTING SERVER ON PORT", cfg.port)

	err = server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
