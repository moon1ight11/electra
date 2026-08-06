package app

import (
	"electra/internal/api"
	authhandlers "electra/internal/api/handlers/auth_handlers"
	orderhandlers "electra/internal/api/handlers/order_handlers"
	orderworkerhandlers "electra/internal/api/handlers/order_worker_handlers"
	requesthandlers "electra/internal/api/handlers/request_handlers"
	statisticshandlers "electra/internal/api/handlers/statistics_handlers"
	workerhandlers "electra/internal/api/handlers/worker_handlers"
	"electra/internal/api/jwt"
	"electra/internal/config"
	authservice "electra/internal/services/auth_service"
	orderservice "electra/internal/services/order_service"
	orderworkerservice "electra/internal/services/order_worker_service"
	requestservice "electra/internal/services/request_service"
	statisticservice "electra/internal/services/statistic_service"
	"electra/internal/storage/repo/database"
	orderworkers "electra/internal/storage/repo/order_workers"
	"electra/internal/storage/repo/orders"
	"electra/internal/storage/repo/requests"
	"electra/internal/storage/repo/statistic"
	"electra/internal/storage/repo/workers"
	"electra/pkg/logger"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	cfg    config.Config
	logger logger.Logger
	db     *database.DataBase
	router *api.Router
}

func New(cfg config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Init() error {
	var err error

	a.logger, err = logger.NewLogger(a.cfg)
	if err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}

	a.db, err = database.PostgresConnection(a.cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := a.db.UpMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	a.logger.Info("migrations applied")

	jwtService := jwt.NewJWTService(
		a.cfg.JWT.Secret,
		time.Duration(a.cfg.JWT.Expiration)*time.Hour,
	)

	requestRepo := requests.NewRequestRepo(a.db)
	workerRepo := workers.NewWorkerRepo(a.db)
	orderRepo := orders.NewOrderRepo(a.db)
	orderWorkerRepo := orderworkers.NewOrderWorkerRepo(a.db)
	statisticsRepo := statistic.NewStatisticsRepo(a.db)

	authService := authservice.NewAuthService(workerRepo, jwtService)
	requestService := requestservice.NewRequestService(requestRepo, orderRepo)
	orderService := orderservice.NewOrderService(orderRepo, orderWorkerRepo)
	orderWorkerService := orderworkerservice.NewOrderWorkerService(orderWorkerRepo, workerRepo)
	statisticsService := statisticservice.NewStatisticsService(statisticsRepo)

	authHandler := authhandlers.NewAuthHandler(authService, a.logger)
	requestHandler := requesthandlers.NewRequestHandler(requestService, a.logger)
	orderHandler := orderhandlers.NewOrderHandler(orderService, a.logger)
	orderWorkerHandler := orderworkerhandlers.NewOrderWorkerHandler(orderWorkerService, a.logger)
	statisticsHandler := statisticshandlers.NewStatisticsHandler(statisticsService, a.logger)
	workerHandler := workerhandlers.NewWorkerHandler(authService, a.logger)

	a.router = api.NewRouter(
		authHandler,
		requestHandler,
		orderHandler,
		orderWorkerHandler,
		statisticsHandler,
		workerHandler,
	)
	a.router.Init(jwtService, a.logger, a.cfg.Server.CORSOrigin)

	return nil
}

func (a *App) Run() error {
	addr := fmt.Sprintf("%s:%s", a.cfg.Server.Host, a.cfg.Server.Port)
	a.logger.Info("starting server", slog.String("addr", addr))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := a.router.GetEngine().Run(addr); err != nil {
			a.logger.Error("server error", slog.String("error", err.Error()))
		}
	}()

	<-quit
	a.logger.Info("shutting down server...")

	return a.db.Close()
}