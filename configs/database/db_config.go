package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type DBConfig struct {
	DB_addr      string // e.g., "file:app.db?cache=shared&_pragma=busy_timeout=5000&_pragma=journal_mode=WAL"
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string // e.g., "30m", "1h"
}

// NewSQLiteDB initializes a Bun DB instance with SQLite and connection pooling
func NewSQLiteDB(cfg DBConfig, dropFlag bool) (*bun.DB, error) {
	// 1. Open standard sql.DB using Bun's recommended pure-Go SQLite shim (no CGO required)
	sqldb, err := sql.Open(sqliteshim.DriverName(), cfg.DB_addr)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// 2. Configure connection pool settings
	// Note: For SQLite, MaxOpenConns is often set to 1 for write-heavy loads 
	// to prevent disk contention, but can be scaled for read-heavy WAL mode.
	sqldb.SetMaxOpenConns(cfg.MaxOpenConns)
	sqldb.SetMaxIdleConns(cfg.MaxIdleConns)

	if cfg.MaxIdleTime != "" {
		idleDuration, err := time.ParseDuration(cfg.MaxIdleTime)
		if err != nil {
			sqldb.Close()
			return nil, fmt.Errorf("invalid MaxIdleTime format: %w", err)
		}
		sqldb.SetConnMaxIdleTime(idleDuration)
	}

	// 3. Wrap standard sql.DB with Bun using the SQLite dialect
	db := bun.NewDB(sqldb, sqlitedialect.New())

	// 4. Verify connection with a 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	// 5. Handle database initialization and drop flag schema reset
	if err := InitDatabase(db, dropFlag); err != nil {
		db.Close()
		return nil, fmt.Errorf("database initialization failed: %w", err)
	}

	slog.Info("Connected to SQLite via Bun!",
		"addr", cfg.DB_addr,
		"max_open_conns", cfg.MaxOpenConns,
		"max_idle_conns", cfg.MaxIdleConns,
	)

	return db, nil
}