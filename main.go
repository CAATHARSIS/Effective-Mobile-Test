package main

import (
	"Effective-Mobile-Test/internal/config"
	"Effective-Mobile-Test/pkg/database"
	"fmt"
	"log/slog"
	"os"

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
}
