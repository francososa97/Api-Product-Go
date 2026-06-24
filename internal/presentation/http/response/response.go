// Package response centraliza la escritura de respuestas HTTP en formato JSON,
// garantizando un Content-Type correcto y un formato de error consistente en
// toda la API.
package response

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/francososa97/product-api/internal/domain"
)

// ErrorBody es el cuerpo estándar de toda respuesta de error de la API.
type ErrorBody struct {
	Error      string `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

// JSON escribe v como JSON con el status indicado.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// La cabecera ya fue enviada; solo queda dejar traza del problema.
		log.Printf("error al codificar respuesta JSON: %v", err)
	}
}

// Error escribe una respuesta de error con cuerpo consistente.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorBody{
		Error:      http.StatusText(status),
		Message:    message,
		StatusCode: status,
	})
}

// FromDomainError traduce un error de dominio al código HTTP que le corresponde,
// manteniendo al dominio totalmente ignorante del transporte.
func FromDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrEmptyName), errors.Is(err, domain.ErrNegativePrice):
		Error(w, http.StatusBadRequest, err.Error())
	default:
		// No se filtran detalles internos al cliente; el detalle va al log.
		log.Printf("error interno: %v", err)
		Error(w, http.StatusInternalServerError, "error interno del servidor")
	}
}
