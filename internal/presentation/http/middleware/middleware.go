// Package middleware agrupa middlewares HTTP transversales (logging, recover).
package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/francososa97/product-api/internal/presentation/http/response"
)

// statusRecorder envuelve un ResponseWriter para capturar el código de estado
// efectivamente escrito, necesario para el logging.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Logging registra método, ruta, código de estado y duración de cada request.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, rec.status, time.Since(start))
	})
}

// Recover captura panics en los handlers y responde 500 en lugar de tirar la
// conexión, dejando traza del incidente en el log.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic recuperado en %s %s: %v", r.Method, r.URL.Path, rec)
				response.Error(w, http.StatusInternalServerError, "error interno del servidor")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
