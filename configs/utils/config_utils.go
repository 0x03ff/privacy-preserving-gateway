package configutils

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	fame "github.com/0x03ff/gateway/third_party"
)

// SetupLogging configures structured logging using log/slog, writing to both stdout and a log file.
func SetupLogging() (*slog.Logger, *os.File, error) {
	// Create logs directory if missing
	if err := os.MkdirAll("./configs/logs", 0755); err != nil {
		return nil, nil, err
	}

	// Open log file with append mode
	logFile, err := os.OpenFile("./configs/logs/system.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, err
	}

	// Create multi-writer (terminal + log file)
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Configure handler options
	opts := &slog.HandlerOptions{
		AddSource: true,             // Equivalent to log.Lshortfile (includes file and line number)
		Level:     slog.LevelInfo,   // Minimum log level (Debug, Info, Warn, Error)
	}

	// Choose TextHandler for human-readable output (or JSONHandler for production)
	handler := slog.NewTextHandler(multiWriter, opts)

	// Create the logger instance
	logger := slog.New(handler)

	// Set as the global default logger so package-level slog calls work anywhere
	slog.SetDefault(logger)

	// Return the custom logger instance, log file handle, and nil error
	return logger, logFile, nil
}

func SetupFame() error {
	dir := "configs/fame"
	fame_dir := "assets/fame"
	files := []string{"fame.gob", "publicKey.gob", "secretKey.gob"}

	// Create directories if missing
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(fame_dir, 0755); err != nil {
		return err
	}

	// Check if all required files exist
	missing := false
	for _, file := range files {
		if _, err := os.Stat(filepath.Join(dir, file)); os.IsNotExist(err) {
			missing = true
			break
		} else if err != nil {
			return err
		}
	}

	// Generate keys if any file is missing using structured slog
	if missing {
		slog.Info("FAME setup: Generating new master keys (missing files detected)")
		if err := fame.GenerateMasterKeys(dir, true); err != nil {
			return err
		}
		slog.Info("FAME setup: Master keys generated successfully")
	} else {
		slog.Info("FAME setup: All required files exist - no action needed")
	}

	// Copy public files to assets directory
	publicFiles := []string{"fame.gob", "publicKey.gob"}
	for _, file := range publicFiles {
		srcPath := filepath.Join(dir, file)
		dstPath := filepath.Join(fame_dir, file)

		// Check if source file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			return fmt.Errorf("source file missing: %s", srcPath)
		} else if err != nil {
			return err
		}

		// Copy file to assets
		if err := copyFile(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to copy %s to %s: %w", srcPath, dstPath, err)
		}
		
		// Log file copy with structured attributes
		slog.Info("FAME setup: Copied file successfully", "src", srcPath, "dst", dstPath)
	}

	return nil
}

func GetFamePublicFiles() (fameData []byte, publicKeyData []byte, er error) {
	// Read fame.gob
	fame, err := os.ReadFile("configs/fame/fame.gob")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read fame.gob: %w", err)
	}

	// Read publicKey.gob
	pubKey, err := os.ReadFile("configs/fame/publicKey.gob")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read publicKey.gob: %w", err)
	}

	return fame, pubKey,
		nil
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


