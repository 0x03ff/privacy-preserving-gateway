package configutils

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// provide default value if the varible not available

func GetEnvString(key, fallback string) string {

	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

func GetEnvInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return valAsInt
}

// Critical variable validation
type MissingEnvError struct {
	Key string
}

func (e *MissingEnvError) Error() string {
	msg := fmt.Sprintf("missing required environment variable: %s", e.Key)
	log.Println(msg)
	return msg
}

// MustGetString returns error if variable is missing (no fallback)
func MustGetEnvString(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", &MissingEnvError{Key: key}
	}
	return val, nil
}

// MustGetInt returns error if variable is missing or invalid
func MustGetEnvInt(key string) (int, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return 0, &MissingEnvError{Key: key}
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid integer for %s: %w", key, err)
	}
	return valAsInt, nil
}

// IsMissingEnv checks if error is missing environment variable
func IsMissingEnv(err error) bool {
	_, ok := err.(*MissingEnvError)
	return ok
}
