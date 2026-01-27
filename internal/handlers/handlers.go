package handlers

import (
	"Effective-Mobile-Test/internal/models"
	"Effective-Mobile-Test/internal/repository"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type SubscriptionHandler struct {
	repo repository.RepositoryInterface
	log  *slog.Logger
}

func NewSubscriptionHandler(repo repository.RepositoryInterface, log *slog.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		repo: repo,
		log:  log,
	}
}

func (h *SubscriptionHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/subscriptions", h.CreateSubscriptionRecord).Methods("POST")
}

func (h *SubscriptionHandler) CreateSubscriptionRecord(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5 * time.Second)
	defer cancel()

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	defer r.Body.Close()

	var createSubscription models.CreateSubscriptionRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.handleError(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &createSubscription)
	if err != nil {
		h.handleError(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	subscription, err := createSubscription.ToSubscription()
	if err != nil {
		h.handleError(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	err = h.repo.Create(ctx, subscription)
	if err != nil {
		h.handleError(w, err.Error(), err, http.StatusInternalServerError)
		return
	}

	h.log.Info("Subscription record created successfully", "ID", subscription.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	data, err := json.Marshal(subscription)
	if err != nil {
		h.log.Error("Failed to marshal subscription record", "error", err)
		return
	}
	w.Write(data)
}


func (h *SubscriptionHandler) handleError(w http.ResponseWriter, message string, err error, status int) {
	http.Error(w, message, status)
	h.log.Error(err.Error())
}