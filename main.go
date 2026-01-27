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

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

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

	server := &http.Server{
		Addr: ":" + cfg.ServerPort,
		Handler: router,
	}
	
	server.ListenAndServe()
	defer server.Close()
}
