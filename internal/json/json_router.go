package json

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// JsonHandlers holds all JSON API handlers
type JsonHandlers struct {
	validator *validator.Validate
}

func NewHandlers(validator *validator.Validate) *JsonHandlers {
	return &JsonHandlers{validator: validator}
}

type App interface {
	GetJsonHandlers() *JsonHandlers
}

// SetupJsonRoutes configures all JSON API routes using standard net/http ServeMux
func SetupJsonRoutes(mux *http.ServeMux, app App) {
	handlers := app.GetJsonHandlers()

	mux.HandleFunc("POST /api/login", handlers.AdminLoginHandler)
}

func (h *JsonHandlers) SendError(w http.ResponseWriter, statusCode int, action string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":   statusCode,
		"action": action,
		"error":  message,
	})
}