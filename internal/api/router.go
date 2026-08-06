package api

import (
	authhandlers "electra/internal/api/handlers/auth_handlers"
	healthhandler "electra/internal/api/handlers/health_handler"
	orderhandlers "electra/internal/api/handlers/order_handlers"
	orderworkerhandlers "electra/internal/api/handlers/order_worker_handlers"
	requesthandlers "electra/internal/api/handlers/request_handlers"
	statisticshandlers "electra/internal/api/handlers/statistics_handlers"
	workerhandlers "electra/internal/api/handlers/worker_handlers"
	"electra/internal/api/jwt"
	"electra/internal/api/middlewares"
	"electra/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler        *authhandlers.AuthHandler
	requestHandler     *requesthandlers.RequestHandler
	orderHandler       *orderhandlers.OrderHandler
	orderWorkerHandler *orderworkerhandlers.OrderWorkerHandler
	statisticsHandler  *statisticshandlers.StatisticsHandler
	workerHandler      *workerhandlers.WorkerHandler
	ginEngine          *gin.Engine
}

func NewRouter(
	authHandler *authhandlers.AuthHandler,
	requestHandler *requesthandlers.RequestHandler,
	orderHandler *orderhandlers.OrderHandler,
	orderWorkerHandler *orderworkerhandlers.OrderWorkerHandler,
	statisticsHandler *statisticshandlers.StatisticsHandler,
	workerHandler *workerhandlers.WorkerHandler,
) *Router {
	return &Router{
		authHandler:        authHandler,
		requestHandler:     requestHandler,
		orderHandler:       orderHandler,
		orderWorkerHandler: orderWorkerHandler,
		statisticsHandler:  statisticsHandler,
		workerHandler:      workerHandler,
		ginEngine:          gin.Default(),
	}
}

func (r *Router) Init(jwtService jwt.TokenService, logger logger.Logger, corsOrigin string) {
	r.ginEngine.Use(middlewares.CORS(corsOrigin))

	// группировка роутов
	publicGroup := r.ginEngine.Group("/api/v1/public")
	authGroup := r.ginEngine.Group("/api/v1/auth")
	workerGroup := r.ginEngine.Group("/api/v1/worker")
	ownerGroup := r.ginEngine.Group("/api/v1/owner")

	// мидлвары для аутентификации
	workerGroup.Use(middlewares.Auth(jwtService, logger))
	ownerGroup.Use(middlewares.Auth(jwtService, logger))
	ownerGroup.Use(middlewares.AuthOwner(logger))

	rateLimiter := middlewares.NewRateLimiter(5, time.Minute)
	healthHandler := healthhandler.NewHealthHandler()

	publicGroup.GET("/health", healthHandler.Check)

	// АУТЕНТИФИКАЦИЯ // --- только для созданных пользователей
	authGroup.POST("/login", rateLimiter.Middleware(), r.authHandler.Login)

	// РАБОТНИК // --- только для работников
	workerGroup.POST("/logout", r.authHandler.Logout)
	workerGroup.GET("/orders/planned", r.orderHandler.ListPlanned)
	workerGroup.GET("/orders/history", r.orderWorkerHandler.ListCompleted)
	workerGroup.GET("/orders/:id/reports", r.orderWorkerHandler.GetByOrder)
	workerGroup.PATCH("/orders/report", r.orderWorkerHandler.UpdateReport)
	workerGroup.PATCH("/orders/:id/complete", r.orderHandler.Complete)
	workerGroup.GET("/workers", r.workerHandler.ListWorkers)
	workerGroup.GET("/me", r.workerHandler.GetMe)
	workerGroup.PATCH("/me", r.workerHandler.UpdateProfile)

	// ВЛАДЕЛЕЦ // --- только владелец
	ownerGroup.POST("/workers", r.workerHandler.CreateWorker)
	ownerGroup.GET("/requests/new", r.requestHandler.ListNewRequests)
	ownerGroup.GET("/requests/all", r.requestHandler.ListAllRequests)
	ownerGroup.POST("/requests/convert", r.requestHandler.ConvertToOrder)
	ownerGroup.PATCH("/requests/:id/cancel", r.requestHandler.CancelRequest)
	ownerGroup.POST("/orders/direct", r.orderHandler.CreateDirect)
	ownerGroup.GET("/orders/planned", r.orderHandler.ListAllPlanned)
	ownerGroup.PATCH("/orders/:id/complete", r.orderHandler.CompleteByOwner)
	ownerGroup.GET("/orders/history", r.orderWorkerHandler.ListAllCompleted)
	ownerGroup.DELETE("/orders/:orderId/workers/:workerId", r.orderWorkerHandler.RemoveWorker)
	ownerGroup.PATCH("/orders/edit", r.orderHandler.Update)
	ownerGroup.GET("/statistics/workers/:workerId", r.statisticsHandler.ByWorker)
	ownerGroup.GET("/statistics/all", r.statisticsHandler.AllWorkers)
	ownerGroup.GET("/statistics/summary", r.statisticsHandler.Summary)
	ownerGroup.DELETE("/workers/:id", r.workerHandler.DeleteWorker)

	// ПУБЛИЧНОЕ // ---
	publicGroup.POST("/requests", r.requestHandler.CreateRequest)
}

func (r *Router) GetEngine() *gin.Engine {
	return r.ginEngine
}
