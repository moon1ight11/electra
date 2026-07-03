package requestservice

import (
	"electra/internal/storage/repo/orders"
	"electra/internal/storage/repo/requests"
)

type RequestService struct {
	requestRepo *requests.RequestRepo
	orderRepo   *orders.OrderRepo
}

func NewRequestService(requestRepo *requests.RequestRepo, orderRepo *orders.OrderRepo) *RequestService {
	return &RequestService{
		requestRepo: requestRepo,
		orderRepo:   orderRepo,
	}
}
