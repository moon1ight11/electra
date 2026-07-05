package workers

import (
	"context"
	"database/sql"
	"electra/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

// создание исполнителя
func (r *WorkerRepo) Create(ctx context.Context, w *domain.Worker) error {
	query := `
			INSERT INTO workers (name, phone, role, password_hash)
		 	VALUES ($1, $2, $3, $4)
		 	RETURNING id, created_at
			`

	err := r.db.DB.QueryRowContext(ctx, query, w.Name, w.Phone, w.Role, w.PasswordHash).Scan(&w.ID, &w.CreatedAt)
	if err != nil {
		return fmt.Errorf("error in create worker: %w", err)
	}

	return nil
}

// получение исполнителя по id
func (r *WorkerRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Worker, error) {
	w := &domain.Worker{}

	query := `
			SELECT id, name, phone, role, password_hash, created_at
		 	FROM workers
		 	WHERE id = $1
			`

	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(&w.ID, &w.Name, &w.Phone, &w.Role, &w.PasswordHash, &w.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error in get worker: %w", err)
	}

	return w, nil
}

// поиск исполнителя по телефону (для входа)
func (r *WorkerRepo) GetByPhone(ctx context.Context, phone string) (*domain.Worker, error) {
	w := &domain.Worker{}

	query := `
			SELECT id, name, phone, role, password_hash, created_at
		 	FROM workers
		 	WHERE phone = $1
			`

	err := r.db.DB.QueryRowContext(ctx, query, phone).Scan(&w.ID, &w.Name, &w.Phone, &w.Role, &w.PasswordHash, &w.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error in get worker by phone: %w", err)
	}

	return w, nil
}

// список всех исполнителей
func (r *WorkerRepo) List(ctx context.Context) ([]domain.Worker, error) {
	query := `
			SELECT id, name, phone, role, password_hash, created_at
		 	FROM workers
		 	ORDER BY created_at
			`

	rows, err := r.db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error in list workers: %w", err)
	}
	defer rows.Close()

	workers := make([]domain.Worker, 0)
	for rows.Next() {
		var w domain.Worker
		if err := rows.Scan(&w.ID, &w.Name, &w.Phone, &w.Role, &w.PasswordHash, &w.CreatedAt); err != nil {
			return nil, fmt.Errorf("error in scan worker: %w", err)
		}
		workers = append(workers, w)
	}

	return workers, rows.Err()
}
