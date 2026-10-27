package configs

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"time"

	"github.com/0x03ff/gateway/configs/router"
	auth "github.com/0x03ff/gateway/internal/auth"
	json "github.com/0x03ff/gateway/internal/json"
)



type TlsConfig struct {
	CertFile   string
	KeyFile    string
	MinVersion uint16
	TLSConfig  *tls.Config
}

type Application struct {
	Tlsconfig *tls.Config
	Cert_path string
	Key_path  string

	JsonHandlers *json.JsonHandlers
	AuthHandlers *auth.AuthHandlers
}

func (t *TlsConfig) NewTLSConfig() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(t.CertFile, t.KeyFile)
	if err != nil {
		return nil, err
	}

	t.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   t.MinVersion,
	}

	return t.TLSConfig, nil
}

func (app *Application) GetAuthHandlers() *auth.AuthHandlers {
	return app.AuthHandlers
}

func (app *Application) GetJsonHandlers() *json.JsonHandlers {
	return app.JsonHandlers
}

func (app *Application) Mount() http.Handler {
	setupFunc := func(mux *http.ServeMux) {
		auth.SetupAuthRoutes(mux, app)

		// html.SetupHtmlRoutes(mux, app) // Uncomment if html package is used
		json.SetupJsonRoutes(mux, app)
	}

	return router.SetupRoutes(setupFunc)
}

func (app *Application) Run(mux http.Handler, addr string, httpsAddr string) error {
	// HTTP server (port 8086 / pprof)
	httpSrv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		WriteTimeout: time.Minute * 1,
		ReadTimeout:  time.Minute * 1,
		IdleTimeout:  time.Minute,
	}

	go func() {
		slog.Info("HTTP/pprof server started", "addr", addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server failed", "error", err)
		}
	}()

	// HTTPS server (port 443)
	httpsSrv := &http.Server{
		Addr:         httpsAddr,
		Handler:      mux,
		WriteTimeout: time.Minute * 10,
		ReadTimeout:  time.Minute * 10,
		IdleTimeout:  time.Minute,
		TLSConfig:    app.Tlsconfig,
	}

	// Start HTTPS server (blocks main thread)
	slog.Info("HTTPS server started", "addr", httpsAddr)
	return httpsSrv.ListenAndServeTLS(app.Cert_path, app.Key_path)
}