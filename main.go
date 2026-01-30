package main

import (
	"Effective-Mobile-Test/internal/config"
	"Effective-Mobile-Test/internal/handlers"
	"Effective-Mobile-Test/internal/repository"
	"Effective-Mobile-Test/pkg/database"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "Effective-Mobile-Test/docs"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Effective-Mobile-Test API
// @version 1.0
// @description REST API для агрегации данных об онлайн подписках пользователей

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemas http

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	cfg := config.Load(log)

	fmt.Println(cfg)

	migrationDB, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
	}

	if err := database.RunMigrations(migrationDB, log); err != nil {
		log.Error("Failed to run migrations", "error", err)
		if err := migrationDB.Close(); err != nil {
			log.Error("Failed to close migration connection", "error", err)
		}
	}

	appDB, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
	}

	repo := repository.NewSubscriptionRepo(appDB)
	handler := handlers.NewSubscriptionHandler(repo, log)

	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	server.ListenAndServe()
	log.Info("Server started", "port", cfg.ServerPort)

	defer server.Close()
}
