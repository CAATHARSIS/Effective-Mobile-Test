package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int        `json:"int"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type SubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

type SubscriptionResponse struct {
	ID          int     `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

type SubscriptionCostRequest struct {
	ServiceName string `json:"service_name"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type SubscriptionCost struct {
	ServiceName sql.NullString `json:"service_name"`
	UserID      *uuid.UUID     `json:"user_id"`
	StartDate   *time.Time     `json:"start_date"`
	EndDate     *time.Time     `json:"end_date"`
}

type SubscriptionCostResponse struct {
	Cost int `json:"cost"`
}

func (s SubscriptionCost) String() string {
	return fmt.Sprintf("start_date: %v\nend_date: %v", s.StartDate, s.EndDate)
}

func (req SubscriptionCostRequest) ToSubscriptionCost() (*SubscriptionCost, error) {
	var subscription SubscriptionCost

	if req.ServiceName != "" {
		subscription.ServiceName.String = req.ServiceName
		subscription.ServiceName.Valid = true
	}

	if req.UserID != "" {
		userUUID, err := uuid.Parse(req.UserID)
		if err != nil {
			return nil, fmt.Errorf("Invalid user_id format, must be uuid: %v", err)
		}

		subscription.UserID = &userUUID
	}

	if req.StartDate != "" {
		startDate, err := parseDate(req.StartDate)
		if err != nil {
			return nil, err
		}

		if !startDate.IsZero() {
			subscription.StartDate = &startDate
		}
	}

	if req.EndDate != "" {
		endDate, err := parseDate(req.EndDate)
		if err != nil {
			return nil, err
		}

		if !endDate.IsZero() {
			subscription.EndDate = &endDate
		}
	}

	return &subscription, nil
}

func (sub Subscription) ToResponse() *SubscriptionResponse {
	resp := SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID.String(),
		StartDate:   formatDate(sub.StartDate),
	}

	if sub.EndDate != nil {
		temp := formatDate(*sub.EndDate)
		resp.EndDate = &temp
	}

	return &resp
}

func (req SubscriptionRequest) ToSubscription() (*Subscription, error) {
	stardDate, err := parseDate(req.StartDate)
	if err != nil {
		return nil, err
	}

	var endDate time.Time
	if req.EndDate != nil {
		endDate, err = parseDate(*req.EndDate)
		if err != nil {
			return nil, err
		}
	}

	subscription := Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   stardDate,
	}

	if endDate.IsZero() {
		subscription.EndDate = nil
	} else {
		subscription.EndDate = &endDate
	}

	if req.UserID != "" {
		userUUID, err := uuid.Parse(req.UserID)
		if err != nil {
			return nil, fmt.Errorf("Invalid user_id format, must be uuid: %v", err)
		}
		subscription.UserID = userUUID
	}

	return &subscription, nil
}

func parseDate(date string) (time.Time, error) {
	if date == "" {
		return time.Time{}, nil
	}

	parts := strings.Split(date, "-")
	if len(parts) != 2 {
		return time.Time{}, errors.New("Invalid date format, should be like MM-YYYY")
	}

	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, errors.New("Invalid date format, should be like MM-YYYY")
	}
	month, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, errors.New("Invalid date format, should be like MM-YYYY")
	}

	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), nil
}

func formatDate(date time.Time) string {
	return fmt.Sprintf("%02d-%04d", date.Month(), date.Year())
}
