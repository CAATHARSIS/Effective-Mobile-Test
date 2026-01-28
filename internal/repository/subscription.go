package repository

import (
	"Effective-Mobile-Test/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type RepositoryInterface interface {
	Create(ctx context.Context, subscription *models.Subscription) error
	GetByID(ctx context.Context, id int) (*models.Subscription, error)
	Update(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error)
	DeleteByID(ctx context.Context, id int) error
	List(ctx context.Context) ([]*models.Subscription, error)
}

type SubscriptionRepo struct {
	db *sql.DB
}

func NewSubscriptionRepo(db *sql.DB) RepositoryInterface {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Create(ctx context.Context, subscription *models.Subscription) error {
	query := `
		INSERT INTO
			subscription_record (
				service_name,
				price,
				user_id,
				start_date,
				end_date
			)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
	).Scan(&subscription.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *SubscriptionRepo) GetByID(ctx context.Context, id int) (*models.Subscription, error) {
	query := `
		SELECT
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date
		FROM
			subscription_record
		WHERE
			id = $1
	`
	var subscription models.Subscription
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subscription.ID,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserID,
		&subscription.StartDate,
		&subscription.EndDate,
	)

	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
	query := `
		UPDATE subscription_record
		SET
			service_name = $1,
			price = $2,
			user_id = $3,
			start_date = $4,
			end_date = $5
		WHERE id = $6
		RETURNING
			service_name,
			price,
			user_id,
			start_date,
			end_date
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
		subscription.ID,
	).Scan(
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserID,
		&subscription.StartDate,
		&subscription.EndDate,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Subscription record with id %d not found", subscription.ID)
		} else {
			return nil, err
		}
	}

	return subscription, nil
}

func (r *SubscriptionRepo) DeleteByID(ctx context.Context, id int) error {
	query := `
		DELETE FROM subscription_record
		WHERE id = $1
	`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("No subscription record with id: %v", err)
	}

	return nil
}

func (r *SubscriptionRepo) List(ctx context.Context) ([]*models.Subscription, error) {
	query := `
		SELECT
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date
		FROM
			subscription_record
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*models.Subscription
	for rows.Next() {
		var record models.Subscription

		err := rows.Scan(
			&record.ID,
			&record.ServiceName,
			&record.Price,
			&record.UserID,
			&record.StartDate,
			&record.EndDate,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan subscription record while listing: %v", err)
		}

		records = append(records, &record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error while listing: %v", err)
	}

	return records, nil
}
