package orderworkerservice

import (
	orderworkers "electra/internal/storage/repo/order_workers"
	"electra/internal/storage/repo/workers"
)

type OrderWorkerService struct {
	orderWorkerRepo *orderworkers.OrderWorkerRepo
	workerRepo      *workers.WorkerRepo
}

func NewOrderWorkerService(orderWorkerRepo *orderworkers.OrderWorkerRepo, workerRepo *workers.WorkerRepo) *OrderWorkerService {
	return &OrderWorkerService{
		orderWorkerRepo: orderWorkerRepo,
		workerRepo:      workerRepo,
	}
}
