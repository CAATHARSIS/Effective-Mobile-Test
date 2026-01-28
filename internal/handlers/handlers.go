package handlers

import (
	"Effective-Mobile-Test/internal/models"
	"Effective-Mobile-Test/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
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
	router.HandleFunc("/subscriptions/{id}", h.GetSubscriptionRecord).Methods("GET")
	router.HandleFunc("/subscriptions", h.ListSubsriptionRecords).Methods("GET")
	router.HandleFunc("/subscriptions/{id}", h.UpdateSubscriptionRecord).Methods("PUT")
}

func (h *SubscriptionHandler) CreateSubscriptionRecord(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	defer r.Body.Close()

	var createSubscription models.SubscriptionRequest

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
		h.handleError(w, "Failed to create subscription record", err, http.StatusInternalServerError)
		return
	}

	h.log.Info("Subscription record created successfully", "ID", subscription.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	subscriptionResponse := subscription.ToResponse()

	data, err := json.Marshal(subscriptionResponse)
	if err != nil {
		h.log.Error("Failed to marshal subscription record", "error", err)
		return
	}
	w.Write(data)
}

func (h *SubscriptionHandler) GetSubscriptionRecord(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	defer r.Body.Close()

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.handleError(w, "Invalid id in request", err, http.StatusBadRequest)
		return
	}

	subscription, err := h.repo.GetByID(ctx, id)
	if err != nil {
		h.handleError(w, fmt.Sprintf("Failed to get subscription record with id: %d", id), err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	subscriptionResponse := subscription.ToResponse()
	data, err := json.Marshal(subscriptionResponse)
	if err != nil {
		h.handleError(w, "Failed to marshal response", err, http.StatusInternalServerError)
		return
	}
	w.Write(data)
	h.log.Info("Subscription record got successfully", "ID", subscriptionResponse.ID)
}

func (h *SubscriptionHandler) ListSubsriptionRecords(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	defer r.Body.Close()

	subscriptions, err := h.repo.List(ctx)
	if err != nil {
		h.handleError(w, "Failed to list subsription records", err, http.StatusInternalServerError)
		return
	}

	var response []*models.SubscriptionResponse
	for _, sub := range subscriptions {
		curr := sub.ToResponse()
		response = append(response, curr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	h.log.Info("Subscription records listed successfully", "amount", len(response))
}

func (h *SubscriptionHandler) UpdateSubscriptionRecord(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	defer r.Body.Close()

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.handleError(w, "Invalid id in request", err, http.StatusBadRequest)
		return
	}

	var updateSubscription models.SubscriptionRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.handleError(w, "Invalid request", err, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &updateSubscription)
	if err != nil {
		h.handleError(w, "Ivnalid request", err, http.StatusBadRequest)
		return
	}

	subscription, err := updateSubscription.ToSubscription()
	if err != nil {
		h.handleError(w, "Invalid request", err, http.StatusBadRequest)
		return
	}
	subscription.ID = id

	updatedSubscription, err := h.repo.Update(ctx, subscription)
	if err != nil {
		h.handleError(w, "Failed to update subscription record", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	subscriptionResponse := updatedSubscription.ToResponse()
	data, err := json.Marshal(subscriptionResponse)
	if err != nil {
		h.handleError(w, "Failed to marshal updated subscription record", err, http.StatusInternalServerError)
		return
	}
	w.Write(data)

	h.log.Info("Subscription record updated successfully", "ID", id)
}

func (h *SubscriptionHandler) handleError(w http.ResponseWriter, message string, err error, status int) {
	http.Error(w, message, status)
	h.log.Error(message, "error", err)
}
