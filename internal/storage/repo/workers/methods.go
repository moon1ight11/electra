package workers

import (
	"context"
	"database/sql"
	"electra/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

func (r *WorkerRepo) Create(ctx context.Context, w *domain.Worker) error {
	query := `
		INSERT INTO workers (name, phone, role, password_hash, specialization)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	err := r.db.DB.QueryRowContext(ctx, query, w.Name, w.Phone, w.Role, w.PasswordHash, w.Specialization).
		Scan(&w.ID, &w.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create worker: %w", err)
	}
	return nil
}

func (r *WorkerRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Worker, error) {
	w := &domain.Worker{}

	query := `
		SELECT id, name, phone, role, password_hash, specialization, created_at
		FROM workers
		WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.DB.QueryRowContext(ctx, query, id).
		Scan(&w.ID, &w.Name, &w.Phone, &w.Role, &w.PasswordHash, &w.Specialization, &w.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get worker: %w", err)
	}
	return w, nil
}

func (r *WorkerRepo) GetByPhone(ctx context.Context, phone string) (*domain.Worker, error) {
	w := &domain.Worker{}

	query := `
		SELECT id, name, phone, role, password_hash, specialization, created_at
		FROM workers
		WHERE phone = $1 AND deleted_at IS NULL`

	err := r.db.DB.QueryRowContext(ctx, query, phone).
		Scan(&w.ID, &w.Name, &w.Phone, &w.Role, &w.PasswordHash, &w.Specialization, &w.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get worker by phone: %w", err)
	}
	return w, nil
}

func (r *WorkerRepo) List(ctx context.Context) ([]domain.Worker, error) {
	query := `
		SELECT id, name, phone, role, password_hash, specialization, created_at
		FROM workers
		WHERE deleted_at IS NULL
		ORDER BY created_at`

	rows, err := r.db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list workers: %w", err)
	}
	defer rows.Close()

	workers := make([]domain.Worker, 0)
	for rows.Next() {
		var w domain.Worker
		if err := rows.Scan(&w.ID, &w.Name, &w.Phone, &w.Role, &w.PasswordHash, &w.Specialization, &w.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan worker: %w", err)
		}
		workers = append(workers, w)
	}
	return workers, rows.Err()
}

func (r *WorkerRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE workers SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete worker: %w", err)
	}
	return nil
}

func (r *WorkerRepo) Update(ctx context.Context, id uuid.UUID, name string, specialization *string, passwordHash *string) error {
	query := `
		UPDATE workers
		SET name = $1, specialization = $2, password_hash = COALESCE($3, password_hash)
		WHERE id = $4 AND deleted_at IS NULL`

	_, err := r.db.DB.ExecContext(ctx, query, name, specialization, passwordHash, id)
	if err != nil {
		return fmt.Errorf("failed to update worker: %w", err)
	}
	return nil
}