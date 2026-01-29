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

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
	}

	var s string
	err = db.QueryRow("Select current_database()").Scan(&s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s)

	repo := repository.NewSubscriptionRepo(db)
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
	defer server.Close()
}
