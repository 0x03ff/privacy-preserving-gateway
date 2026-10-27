package main

import (
	"crypto/tls"
	"log/slog"

	_ "net/http/pprof" // Import this for pprof endpoints

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	"github.com/0x03ff/gateway/configs"
	"github.com/0x03ff/gateway/configs/database"
	configutils "github.com/0x03ff/gateway/configs/utils"
	"github.com/0x03ff/gateway/internal/auth"
	"github.com/0x03ff/gateway/internal/json"
)

func main() {
	// 1. Load Environment
	if err := godotenv.Load("configs/env/.env"); err != nil {
		slog.Warn("No .env file found, relying on system environment variables")
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	dropFlag := false // Set to true if you want to wipe/reset the SQLite database on startup

	// 2. Logging & Setup
	_, logFile, err := configutils.SetupLogging()
	if err != nil {
		slog.Error("Failed to setup logging", "error", err)
	}
	defer logFile.Close()

	if err := configutils.SetupFame(); err != nil {
		slog.Error("Failed to setup FAME keys", "error", err)
	}

	// 3. Initialize SQLite Database via Bun ORM
	dbAddr := configutils.GetEnvString("DB_ADDR", "file:configs/database/app.db?cache=shared&_pragma=busy_timeout=5000&_pragma=journal_mode=WAL")
	
	dbConfig := database.DBConfig{
		DB_addr:      dbAddr,
		MaxOpenConns: 50,
		MaxIdleConns: 10,
		MaxIdleTime:  "15m",
	}

	db, err := database.NewSQLiteDB(dbConfig, dropFlag)
	if err != nil {
		slog.Error("Failed to connect to SQLite database", "error", err)
		return
	}
	defer db.Close()

	// 4. Load TLS Configuration
	tlsCfg := &configs.TlsConfig{
		CertFile:   "configs/ssl/fullchain.pem",
		KeyFile:    "configs/ssl/go_key.pem",
		MinVersion: tls.VersionTLS12,
	}
	tls_config, err := tlsCfg.NewTLSConfig()
	if err != nil {
		slog.Error("Failed to load TLS certificates", "error", err)
	}

	// 5. Initialize Application
	app := &configs.Application{
		Tlsconfig:    tls_config,
		Cert_path:    "configs/ssl/fullchain.pem",
		Key_path:     "configs/ssl/go_key.pem",
		JsonHandlers: json.NewHandlers(validate),
		AuthHandlers: auth.NewHandlers(validate),
	}

	// 6. Mount standard net/http router mux
	mux := app.Mount()

	// 7. Get Server Addresses from Environment
	addr := configutils.GetEnvString("HTTP", "0.0.0.0:8086")
	httpsAddr := configutils.GetEnvString("HTTPS", "0.0.0.0:8080")

	// 8. Run Servers
	slog.Info("Initializing gateway application...")
	if err := app.Run(mux, addr, httpsAddr); err != nil {
		slog.Error("Server terminated with error", "error", err)
	}
}