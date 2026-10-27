package router

import (
	"log/slog"
	"net/http"
	"time"
)

// SetupRoutesFunc defines the signature for registering routes on the standard ServeMux
type SetupRoutesFunc func(mux *http.ServeMux)

// SetupRoutes initializes the standard router, middleware chain, and static file serving.
func SetupRoutes(setupRoutes SetupRoutesFunc) http.Handler {
	mux := http.NewServeMux()

	// 1. Mount static files (Go 1.22+ path pattern matching)
	fs := http.FileServer(http.Dir("./internal/app/web/assets"))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", fs))
	slog.Info("Serving static files", "dir", "./internal/app/web/assets")

	// 2. Setup application routes via callback function
	setupRoutes(mux)

	// 3. Wrap mux with standard net/http middleware chain
	var handler http.Handler = mux
	handler = recovererMiddleware(handler)
	handler = loggingMiddleware(handler)
	handler = timeoutMiddleware(handler, 10*time.Minute)

	return handler
}

// --- Standard net/http Middleware Implementations ---

// recovererMiddleware catches panics to prevent server crashes
func recovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Panic recovered in request handler", "error", err, "path", r.URL.Path)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs incoming HTTP requests using log/slog
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Capture response status (optional wrapper can be added if exact status code tracking is required)
		next.ServeHTTP(w, r)

		slog.Info("HTTP request processed",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_ip", r.RemoteAddr,
			"duration", time.Since(start),
		)
	})
}

// timeoutMiddleware enforces a request timeout deadline using http.TimeoutHandler
func timeoutMiddleware(next http.Handler, timeout time.Duration) http.Handler {
	return http.TimeoutHandler(next, timeout, "Gateway timeout: request took too long to complete")
}