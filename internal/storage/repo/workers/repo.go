package workers

import "electra/internal/storage/repo/database"

type WorkerRepo struct {
	db *database.DataBase
}

func NewWorkerRepo(db *database.DataBase) *WorkerRepo {
	return &WorkerRepo{db: db}
}
