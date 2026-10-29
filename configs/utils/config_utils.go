package configutils

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	fame "github.com/0x03ff/gateway/third_party"
)

// SetupLogging configures structured logging using log/slog, writing to both stdout and var/logs/system.log.
func SetupLogging() (*slog.Logger, *os.File, error) {
	// Create var/logs directory if missing
	if err := os.MkdirAll("var/logs", 0755); err != nil {
		return nil, nil, err
	}

	// Open log file with append mode
	logFile, err := os.OpenFile("var/logs/system.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, err
	}

	// Create multi-writer (terminal + log file)
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Configure handler options
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}

	handler := slog.NewTextHandler(multiWriter, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger, logFile, nil
}

// SetupFame initializes FAME master keys and saves them under var/fame.
func SetupFame() error {
	dir := "var/fame"
	files := []string{"fame.gob", "publicKey.gob", "secretKey.gob"}

	// Create var/fame directory if missing
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Check if all required FAME files exist
	missing := false
	for _, file := range files {
		if _, err := os.Stat(filepath.Join(dir, file)); os.IsNotExist(err) {
			missing = true
			break
		} else if err != nil {
			return err
		}
	}

	// Generate keys if any file is missing
	if missing {
		slog.Info("FAME setup: Generating new master keys (missing files detected)")
		if err := fame.GenerateMasterKeys(dir, true); err != nil {
			return err
		}
		slog.Info("FAME setup: Master keys generated successfully")
	} else {
		slog.Info("FAME setup: All required FAME files exist - no action needed")
	}

	return nil
}

// GetFamePublicFiles reads public FAME parameters directly from var/fame.
func GetFamePublicFiles() (fameData []byte, publicKeyData []byte, err error) {
	fameBytes, err := os.ReadFile("var/fame/fame.gob")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read fame.gob: %w", err)
	}

	pubKeyBytes, err := os.ReadFile("var/fame/publicKey.gob")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read publicKey.gob: %w", err)
	}

	return fameBytes, pubKeyBytes, nil
}


// Helper function to copy files
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}


