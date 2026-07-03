package orders

import (
	"context"
	"database/sql"
	"electra/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// создание заказа и назначение исполнителей
func (r *OrderRepo) Create(ctx context.Context, o *domain.Order, workerIDs []uuid.UUID) error {
	tx, err := r.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error in begin create order tx: %w", err)
	}
	defer tx.Rollback()

	// создаем заказ
	query := `
			INSERT INTO orders (request_id, address, description, estimated_price, planned_date, created_by)
		 	VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, created_at
			`
	err = tx.QueryRowContext(
		ctx,
		query,
		o.RequestID,
		o.Address,
		o.Description,
		o.EstimatedPrice,
		o.PlannedDate,
		o.CreatedBy).Scan(&o.ID, &o.CreatedAt)
	if err != nil {
		return fmt.Errorf("error in create order: %w", err)
	}

	// назначаем исполнителей
	query = `
			INSERT INTO order_workers (order_id, worker_id) 
			VALUES ($1, $2)
			`

	for _, wid := range workerIDs {
		if _, err = tx.ExecContext(ctx, query, o.ID, wid); err != nil {
			return fmt.Errorf("error in assign worker %s: %w", wid.String(), err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error in commit create order tx: %w", err)
	}

	return nil
}

// возврат заказа по id
func (r *OrderRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	o := &domain.Order{}

	query := `
			SELECT id, request_id, address, description, estimated_price, planned_date, created_by, created_at, completed_at
		 	FROM orders
			WHERE id = $1
			`
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(&o.ID,
		&o.RequestID,
		&o.Address,
		&o.Description,
		&o.EstimatedPrice,
		&o.PlannedDate,
		&o.CreatedBy,
		&o.CreatedAt,
		&o.CompletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error in get order: %w", err)
	}
	return o, nil
}

// запланированные заказы конкретного исполнителя
func (r *OrderRepo) ListPlannedByWorker(ctx context.Context, workerID uuid.UUID) ([]domain.Order, error) {
	query := `
			SELECT o.id, o.request_id, o.address, o.description, o.estimated_price,
		        o.planned_date, o.created_by, o.created_at, o.completed_at
		 	FROM orders o
		 	JOIN order_workers ow ON ow.order_id = o.id
		 	WHERE ow.worker_id = $1 AND o.completed_at IS NULL
		 	ORDER BY o.planned_date
			`
	rows, err := r.db.DB.QueryContext(ctx, query, workerID)
	if err != nil {
		return nil, fmt.Errorf("error in list planned orders: %w", err)
	}
	defer rows.Close()

	orders := make([]domain.Order, 0)
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID,
			&o.RequestID,
			&o.Address,
			&o.Description,
			&o.EstimatedPrice,
			&o.PlannedDate,
			&o.CreatedBy,
			&o.CreatedAt,
			&o.CompletedAt); err != nil {
			return nil, fmt.Errorf("error in scan order: %w", err)
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

// все запланированные заказы (только владелец)
func (r *OrderRepo) ListAllPlanned(ctx context.Context) ([]domain.Order, error) {
	query := `
			SELECT id, request_id, address, description, estimated_price,
		        planned_date, created_by, created_at, completed_at
		 	FROM orders
		 	WHERE completed_at IS NULL
		 	ORDER BY planned_date
			`
	rows, err := r.db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error in list all planned orders: %w", err)
	}
	defer rows.Close()

	orders := make([]domain.Order, 0)
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID,
			&o.RequestID,
			&o.Address,
			&o.Description,
			&o.EstimatedPrice,
			&o.PlannedDate,
			&o.CreatedBy,
			&o.CreatedAt,
			&o.CompletedAt); err != nil {
			return nil, fmt.Errorf("error in scan order: %w", err)
		}
		orders = append(orders, o)
	}

	return orders, rows.Err()
}

// завершение заказа
func (r *OrderRepo) Complete(ctx context.Context, orderID uuid.UUID) error {
	query := `
			UPDATE orders 
			SET completed_at = $1 
			WHERE id = $2 AND completed_at IS NULL
			`
	result, err := r.db.DB.ExecContext(ctx, query, time.Now(), orderID)
	if err != nil {
		return fmt.Errorf("error in complete order: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("error in complete order: order %s not found or already completed", orderID.String())
	}
	return nil
}
