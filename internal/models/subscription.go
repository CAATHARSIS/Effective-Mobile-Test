package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Subscription struct {
	ID          int        `json:"int"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      string     `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type CreateSubscriptionRequest struct {
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

func (sub Subscription) ToResponse() *SubscriptionResponse {
	resp := SubscriptionResponse{
		ID: sub.ID,
		ServiceName: sub.ServiceName,
		Price: sub.Price,
		UserID: sub.UserID,
		StartDate: formatDate(sub.StartDate),
	}

	if sub.EndDate != nil {
		temp := formatDate(*sub.EndDate)
		resp.EndDate = &temp
	}

	return &resp
}

func (req CreateSubscriptionRequest) ToSubscription() (*Subscription, error) {
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
		UserID:      req.UserID,
		StartDate:   stardDate,
	}

	if endDate.IsZero() {
		subscription.EndDate = nil
	} else {
		subscription.EndDate = &endDate
	}

	return &subscription, nil
}

func parseDate(date string) (time.Time, error) {
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
