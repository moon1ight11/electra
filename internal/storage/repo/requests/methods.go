package requests

import (
	"context"
	"database/sql"
	"electra/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

// создание заявки
func (r *RequestRepo) Create(ctx context.Context, req *domain.Request) error {
	query := `
			INSERT INTO requests (name, phone, comment)
		 	VALUES ($1, $2, $3)
		 	RETURNING id, status, created_at`

	err := r.db.DB.QueryRowContext(ctx, query, req.Name, req.Phone, req.Comment).
		Scan(&req.ID, &req.Status, &req.CreatedAt)
	if err != nil {
		return fmt.Errorf("error in create request: %w", err)
	}
	return nil
}

// получение заявки по айди
func (r *RequestRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Request, error) {
	req := &domain.Request{}

	query := `
			SELECT id, name, phone, comment, status, created_at
		 	FROM requests
		 	WHERE id = $1`

	err := r.db.DB.QueryRowContext(ctx, query, id).
		Scan(&req.ID, &req.Name, &req.Phone, &req.Comment, &req.Status, &req.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error in get request by id: %w", err)
	}
	return req, nil
}

// все заявки от новых к старым
func (r *RequestRepo) List(ctx context.Context) ([]domain.Request, error) {
	query := `
			SELECT id, name, phone, comment, status, created_at
		 	FROM requests
		 	ORDER BY created_at DESC`

	rows, err := r.db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error in list requests: %w", err)
	}
	defer rows.Close()

	requests := make([]domain.Request, 0)
	for rows.Next() {
		var req domain.Request
		if err := rows.Scan(&req.ID, &req.Name, &req.Phone, &req.Comment, &req.Status, &req.CreatedAt); err != nil {
			return nil, fmt.Errorf("error in scan request: %w", err)
		}
		requests = append(requests, req)
	}
	return requests, rows.Err()
}

// только новые заявки (для владельца — входящие)
func (r *RequestRepo) ListNew(ctx context.Context) ([]domain.Request, error) {
	query := `
			SELECT id, name, phone, comment, status, created_at
		 	FROM requests
		 	WHERE status = 'new'
		 	ORDER BY created_at DESC`

	rows, err := r.db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error in list new requests: %w", err)
	}
	defer rows.Close()

	requests := make([]domain.Request, 0)
	for rows.Next() {
		var req domain.Request
		if err := rows.Scan(&req.ID, &req.Name, &req.Phone, &req.Comment, &req.Status, &req.CreatedAt); err != nil {
			return nil, fmt.Errorf("error in scan request: %w", err)
		}
		requests = append(requests, req)
	}
	return requests, rows.Err()
}

// пометить заявку как преобразованную в заказ
func (r *RequestRepo) MarkConverted(ctx context.Context, id uuid.UUID) error {
	query := `
			UPDATE requests
			SET status = 'converted'
			WHERE id = $1 AND status = 'new'`

	result, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error in mark request converted: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("request %s not found or not in 'new' status", id.String())
	}
	return nil
}

// пометить заявку как отменённую (клиент отказался)
func (r *RequestRepo) MarkCancelled(ctx context.Context, id uuid.UUID) error {
	query := `
			UPDATE requests
			SET status = 'cancelled'
			WHERE id = $1 AND status = 'new'`

	result, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error in mark request cancelled: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("request %s not found or not in 'new' status", id.String())
	}
	return nil
}
