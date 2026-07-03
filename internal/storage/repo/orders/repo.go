package orders

import "electra/internal/storage/repo/database"

type OrderRepo struct {
	db *database.DataBase
}

// NewOrderRepo — конструктор репозитория заказов.
func NewOrderRepo(db *database.DataBase) *OrderRepo {
	return &OrderRepo{db: db}
}
