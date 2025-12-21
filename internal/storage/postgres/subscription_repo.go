package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DevSchmied/subscription-aggregation-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SubscriptionRepo provides subscription persistence.
type SubscriptionRepo struct {
	pool *pgxpool.Pool
}

// NewSubscriptionRepo creates a new repository instance.
func NewSubscriptionRepo(pool *pgxpool.Pool) *SubscriptionRepo {
	return &SubscriptionRepo{pool: pool}
}

// Create inserts a new subscription.
func (r *SubscriptionRepo) Create(
	ctx context.Context,
	s domain.Subscription,
) (domain.Subscription, error) {

	const q = `
		INSERT INTO subscriptions (
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at;
	`

	// Execute insert and scan timestamps
	if err := r.pool.QueryRow(
		ctx,
		q,
		s.ID,
		s.ServiceName,
		s.Price,
		s.UserID,
		s.StartDate,
		s.EndDate,
	).Scan(&s.CreatedAt, &s.UpdatedAt); err != nil {
		return domain.Subscription{}, fmt.Errorf("create subscription: %w", err)
	}

	return s, nil
}

// GetByID returns subscription by ID.
func (r *SubscriptionRepo) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (domain.Subscription, error) {

	const q = `
		SELECT
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date,
			created_at,
			updated_at
		FROM subscriptions
		WHERE id = $1;
	`

	var s domain.Subscription

	// Query single row by ID
	if err := r.pool.QueryRow(ctx, q, id).Scan(
		&s.ID,
		&s.ServiceName,
		&s.Price,
		&s.UserID,
		&s.StartDate,
		&s.EndDate,
		&s.CreatedAt,
		&s.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subscription{}, ErrNotFound
		}
		return domain.Subscription{}, fmt.Errorf("get subscription by id: %w", err)
	}

	return s, nil
}

// Update modifies an existing subscription.
func (r *SubscriptionRepo) Update(
	ctx context.Context,
	s domain.Subscription,
) (domain.Subscription, error) {

	const q = `
		UPDATE subscriptions
		SET
			service_name = $2,
			price = $3,
			user_id = $4,
			start_date = $5,
			end_date = $6,
			updated_at = now()
		WHERE id = $1
		RETURNING created_at, updated_at;
	`

	// Update fields and timestamps
	if err := r.pool.QueryRow(
		ctx,
		q,
		s.ID,
		s.ServiceName,
		s.Price,
		s.UserID,
		s.StartDate,
		s.EndDate,
	).Scan(&s.CreatedAt, &s.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subscription{}, ErrNotFound
		}
		return domain.Subscription{}, fmt.Errorf("update subscription: %w", err)
	}

	return s, nil
}

// Delete removes subscription by ID.
func (r *SubscriptionRepo) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	const q = `
		DELETE FROM subscriptions
		WHERE id = $1;
	`

	// Execute delete statement
	affectedRows, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}

	// Check affected rows
	if affectedRows.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// ListFilter defines optional list filters.
type ListFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
}

// List returns subscriptions with optional filters.
func (r *SubscriptionRepo) List(
	ctx context.Context,
	f ListFilter,
) ([]domain.Subscription, error) {

	const q = `
		SELECT
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date,
			created_at,
			updated_at
		FROM subscriptions
		WHERE ($1::uuid IS NULL OR user_id = $1)
		  AND ($2::text IS NULL OR service_name = $2)
		ORDER BY created_at DESC;
	`

	// Query filtered subscriptions
	rows, err := r.pool.Query(ctx, q, f.UserID, f.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}
	defer rows.Close()

	var out []domain.Subscription
	for rows.Next() {
		var s domain.Subscription
		if err := rows.Scan(
			&s.ID,
			&s.ServiceName,
			&s.Price,
			&s.UserID,
			&s.StartDate,
			&s.EndDate,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("list scan: %w", err)
		}
		out = append(out, s)
	}

	// Check iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list rows: %w", err)
	}

	return out, nil
}

// ListOverlapping returns subscriptions overlapping period.
func (r *SubscriptionRepo) ListOverlapping(
	ctx context.Context,
	userID *uuid.UUID,
	serviceName *string,
	periodStart,
	periodEnd time.Time,
) ([]domain.Subscription, error) {

	const q = `
		SELECT
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date,
			created_at,
			updated_at
		FROM subscriptions
		WHERE ($1::uuid IS NULL OR user_id = $1)
		  AND ($2::text IS NULL OR service_name = $2)
		  AND start_date <= $4
		  AND (end_date IS NULL OR end_date >= $3)
		ORDER BY start_date ASC;
	`

	// Query overlapping subscriptions
	rows, err := r.pool.Query(ctx, q, userID, serviceName, periodStart, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("list overlapping: %w", err)
	}
	defer rows.Close()

	var out []domain.Subscription
	for rows.Next() {
		var s domain.Subscription
		if err := rows.Scan(
			&s.ID,
			&s.ServiceName,
			&s.Price,
			&s.UserID,
			&s.StartDate,
			&s.EndDate,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("overlapping scan: %w", err)
		}
		out = append(out, s)
	}

	// Check iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("overlapping rows: %w", err)
	}

	return out, nil
}
