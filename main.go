package main

import (
	"Effective-Mobile-Test/internal/config"
	"Effective-Mobile-Test/internal/models"
	"Effective-Mobile-Test/internal/repository"
	"Effective-Mobile-Test/pkg/database"
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
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
	now := time.Now()

	pas := uuid.New().String()
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
	}
	test := &models.Subscription{
		ServiceName: "fdf",
		Price:       32,
		UserID:      pas,
		StartDate:   now,
	}

	err = repo.Create(context.Background(), test)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
	}

	res, err := repo.GetByID(context.Background(), test.ID)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
	}
	log.Info(fmt.Sprintf("%v", *res))

	updatedRes, err := repo.Update(context.Background(), &models.Subscription{
		ID: test.ID,
		ServiceName: "fdf",
		Price:       42,
		UserID:      test.UserID,
		StartDate:   now,
	})

	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
	}
	log.Info(fmt.Sprintf("%v", *updatedRes))

	test2 := &models.Subscription{
		ServiceName: "jhgjg",
		Price:       555,
		UserID:      uuid.New().String(),
		StartDate:   now,
	}
	repo.Create(context.Background(), test2)

	recordsList, err := repo.List(context.Background())
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
	}

	fmt.Println("----------------------------")
	for _, v := range recordsList {
		fmt.Printf("%v\n", *v)
	}

	err = repo.DeleteByID(context.Background(), test.ID)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
	}

	repo.DeleteByID(context.Background(), test2.ID)
}
