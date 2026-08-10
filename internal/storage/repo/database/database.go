package database

import (
	"context"
	"database/sql"
	"electra/internal/config"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type DataBase struct {
	DB            *sql.DB
	MigrationsDir string
}

func (d *DataBase) UpMigrations() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(d.DB, d.MigrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func PostgresConnection(cfg config.Config) (*DataBase, error) {
	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=require",
		cfg.DataBase.User,
		cfg.DataBase.Password,
		cfg.DataBase.Host,
		cfg.DataBase.Port,
		cfg.DataBase.DBName,
	)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(cfg.DataBase.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DataBase.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.DataBase.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.DataBase.ConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DataBase{
		DB:            db,
		MigrationsDir: cfg.DataBase.MigrationsDir,
	}, nil
}

func (d *DataBase) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}
