package requests

import (
	"electra/internal/storage/repo/database"
)

type RequestRepo struct {
	db *database.DataBase
}

func NewRequestRepo(db *database.DataBase) *RequestRepo {
	return &RequestRepo{db: db}
}
