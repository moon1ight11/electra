package database

import (
	"context"
	"database/sql"
	"electra/internal/config"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// структура клиента бд
type DataBase struct {
	DB            *sql.DB
	MigrationsDir string
}

// применение миграций
func (d *DataBase) UpMigrations() error {
	goose.SetBaseFS(nil)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect (upMig): %w", err)
	}

	err := goose.Up(d.DB, d.MigrationsDir)
	if err != nil {
		return fmt.Errorf("failed to up migrations: %w", err)
	}

	return nil
}

// открытие соединения с БД и возврат экземпляра
func PostgresConnection(cfg config.Config) (*DataBase, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DataBase.Host,
		cfg.DataBase.Port,
		cfg.DataBase.User,
		cfg.DataBase.Password,
		cfg.DataBase.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(cfg.DataBase.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DataBase.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.DataBase.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.DataBase.ConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DataBase{
		DB:            db,
		MigrationsDir: cfg.DataBase.MigrationsDir,
	}, nil
}

// закрытие соединения
func (d *DataBase) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}
