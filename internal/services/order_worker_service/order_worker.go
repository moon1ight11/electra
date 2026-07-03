package orderworkerservice

import orderworkers "electra/internal/storage/repo/order_workers"

type OrderWorkerService struct {
	orderWorkerRepo *orderworkers.OrderWorkerRepo
}

func NewOrderWorkerService(orderWorkerRepo *orderworkers.OrderWorkerRepo) *OrderWorkerService {
	return &OrderWorkerService{orderWorkerRepo: orderWorkerRepo}
}
