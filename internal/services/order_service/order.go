package orderservice

import (
	orderworkers "electra/internal/storage/repo/order_workers"
	"electra/internal/storage/repo/orders"
)

type OrderService struct {
	orderRepo       *orders.OrderRepo
	orderWorkerRepo *orderworkers.OrderWorkerRepo
}

func NewOrderService(orderRepo *orders.OrderRepo, orderWorkerRepo *orderworkers.OrderWorkerRepo) *OrderService {
	return &OrderService{
		orderRepo:       orderRepo,
		orderWorkerRepo: orderWorkerRepo,
	}
}
