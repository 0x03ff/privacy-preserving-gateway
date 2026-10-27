package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthHandlers struct {
	validator *validator.Validate

}


func NewHandlers( validator *validator.Validate) *AuthHandlers {
	return &AuthHandlers{ validator: validator}

}


type App interface {
	GetAuthHandlers() *AuthHandlers
}

func SetupAuthRoutes(mux *http.ServeMux, app App) {
	handlers := app.GetAuthHandlers()

	mux.HandleFunc("POST /auth/login", handlers.RegisterHandler)
}

func (h *AuthHandlers) SendError(w http.ResponseWriter, statusCode int, action string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":   statusCode,
		"action": action,
		"error":  message,
	})
}