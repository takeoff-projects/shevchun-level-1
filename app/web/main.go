package main

import (
	"context"
	"embed"
	"github.com/takeoff-projects/level-1/app/web/handlers"
	"github.com/takeoff-projects/level-1/business/data/event"
	"github.com/takeoff-projects/level-1/business/data/schema"
	"github.com/takeoff-projects/level-1/foundation/logger"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func main() {
	log := logger.New()
	defer log.Sync()

	if err := run(log); err != nil {
		log.Error("startup", zap.Error(err))
		os.Exit(1)
	}
}

//go:embed templates/*
var templates embed.FS

func run(log *zap.Logger) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:15434/postgres"
	}
	log.Info("config", zap.String("port", port), zap.String("db_url", dbURL))

	db, err := sqlx.Connect("pgx", dbURL)
	if err != nil {
		log.Error("failed to connect to database", zap.Error(err))
		os.Exit(1)
	}

	if err := schema.Migrate(context.Background(), db); err != nil {
		log.Error("failed to migrate database", zap.Error(err))
		os.Exit(1)
	}

	if err := schema.Seed(context.Background(), db); err != nil {
		log.Error("failed to seed test data", zap.Error(err))
	}

	log.Info("database is up to date")

	fs := http.FileServer(http.Dir("assets"))
	myRouter := mux.NewRouter().StrictSlash(true)

	// This serves the static files in the assets folder
	myRouter.Handle("/assets/", http.StripPrefix("/assets/", fs))

	store := event.NewStore(db)
	h := handlers.Handler{
		Templates: templates,
		Log:       log,
		Store:     store,
	}
	// The rest of the routes
	myRouter.HandleFunc("/", h.IndexHandler)
	myRouter.HandleFunc("/about", h.AboutHandler)
	myRouter.HandleFunc("/add", h.AddHandler)
	myRouter.HandleFunc("/edit/{id}", h.EditHandler)
	myRouter.HandleFunc("/delete/{id}", h.DeleteHandler)

	log.Info("Webserver listening", zap.String("port", port))
	http.ListenAndServe(":"+port, myRouter)

	return nil
}
