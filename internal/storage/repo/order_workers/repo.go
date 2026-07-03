package orderworkers

import "electra/internal/storage/repo/database"

type OrderWorkerRepo struct {
	db *database.DataBase
}

func NewOrderWorkerRepo(db *database.DataBase) *OrderWorkerRepo {
	return &OrderWorkerRepo{db: db}
}
