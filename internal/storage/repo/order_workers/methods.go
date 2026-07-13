package orderworkers

import (
	"context"
	"database/sql"
	"electra/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

// получение списка исполнителей назначенных на заказ
func (r *OrderWorkerRepo) GetByOrder(ctx context.Context, orderID uuid.UUID) ([]domain.OrderWorker, error) {
	query := `
			SELECT id, order_id, worker_id, time_spent, earned_amount, materials_used, notes
		 	FROM order_workers
		 	WHERE order_id = $1
			`

	rows, err := r.db.DB.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("error in get order workers: %w", err)
	}
	defer rows.Close()

	ows := make([]domain.OrderWorker, 0)
	for rows.Next() {
		var ow domain.OrderWorker
		if err := rows.Scan(&ow.ID,
			&ow.OrderID,
			&ow.WorkerID,
			&ow.TimeSpent,
			&ow.EarnedAmount,
			&ow.MaterialsUsed,
			&ow.Notes); err != nil {
			return nil, fmt.Errorf("error in scan order worker: %w", err)
		}
		ows = append(ows, ow)
	}

	return ows, rows.Err()
}

// получение записи заказа исполнителя
func (r *OrderWorkerRepo) GetByOrderAndWorker(ctx context.Context, orderID, workerID uuid.UUID) (*domain.OrderWorker, error) {
	ow := &domain.OrderWorker{}

	query := `
			SELECT id, order_id, worker_id, time_spent, earned_amount, materials_used, notes
		 	FROM order_workers
		 	WHERE order_id = $1 AND worker_id = $2
			`
	err := r.db.DB.QueryRowContext(ctx, query, orderID, workerID).Scan(&ow.ID,
		&ow.OrderID,
		&ow.WorkerID,
		&ow.TimeSpent,
		&ow.EarnedAmount,
		&ow.MaterialsUsed,
		&ow.Notes)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error in get order worker: %w", err)
	}
	return ow, nil
}

// обновление отчета исполнителя по заказу
func (r *OrderWorkerRepo) Update(ctx context.Context, ow *domain.OrderWorker) error {
	query := `
			UPDATE order_workers
		 	SET time_spent = $1, earned_amount = $2, materials_used = $3, notes = $4
		 	WHERE id = $5`
	_, err := r.db.DB.ExecContext(ctx, query, ow.TimeSpent, ow.EarnedAmount, ow.MaterialsUsed, ow.Notes, ow.ID)
	if err != nil {
		return fmt.Errorf("error in update order worker: %w", err)
	}
	return nil
}

// получение всех выполненных заказов исполнителя
func (r *OrderWorkerRepo) ListCompletedByWorker(ctx context.Context, workerID uuid.UUID) ([]domain.Order, error) {
	query := `
		SELECT o.id, o.request_id, o.address, o.description, o.estimated_price,
		       o.planned_date, o.created_by, o.created_at, o.completed_at, rq.phone
		FROM orders o
		JOIN order_workers ow ON ow.order_id = o.id
		LEFT JOIN requests rq ON rq.id = o.request_id
		WHERE ow.worker_id = $1 AND o.completed_at IS NOT NULL
		ORDER BY o.completed_at DESC`

	rows, err := r.db.DB.QueryContext(ctx, query, workerID)
	if err != nil {
		return nil, fmt.Errorf("error in list completed orders: %w", err)
	}
	defer rows.Close()

	orders := make([]domain.Order, 0)
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.RequestID, &o.Address, &o.Description, &o.EstimatedPrice,
			&o.PlannedDate, &o.CreatedBy, &o.CreatedAt, &o.CompletedAt, &o.RequestPhone); err != nil {
			return nil, fmt.Errorf("error in scan order: %w", err)
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

// получение всех выполненных заказов
func (r *OrderWorkerRepo) ListAllCompleted(ctx context.Context) ([]domain.Order, error) {
	query := `
		SELECT o.id, o.request_id, o.address, o.description, o.estimated_price,
		       o.planned_date, o.created_by, o.created_at, o.completed_at, rq.phone
		FROM orders o
		LEFT JOIN requests rq ON rq.id = o.request_id
		WHERE o.completed_at IS NOT NULL
		ORDER BY o.completed_at DESC`

	rows, err := r.db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error in list all completed orders: %w", err)
	}
	defer rows.Close()

	orders := make([]domain.Order, 0)
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.RequestID, &o.Address, &o.Description, &o.EstimatedPrice,
			&o.PlannedDate, &o.CreatedBy, &o.CreatedAt, &o.CompletedAt, &o.RequestPhone); err != nil {
			return nil, fmt.Errorf("error in scan order: %w", err)
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

// удаление исполнителя с заказа
func (r *OrderWorkerRepo) Remove(ctx context.Context, orderID, workerID uuid.UUID) error {
	query := `
			DELETE 
			FROM order_workers 
			WHERE order_id = $1 AND worker_id = $2`

	_, err := r.db.DB.ExecContext(ctx, query, orderID, workerID)
	if err != nil {
		return fmt.Errorf("error in remove order worker: %w", err)
	}
	return nil
}
