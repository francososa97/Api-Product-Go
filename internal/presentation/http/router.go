package http

import (
	"net/http"

	"github.com/francososa97/product-api/internal/application"
	"github.com/francososa97/product-api/internal/presentation/http/docs"
	"github.com/francososa97/product-api/internal/presentation/http/handler"
	"github.com/francososa97/product-api/internal/presentation/http/middleware"
	"github.com/francososa97/product-api/internal/presentation/http/response"
)

// NewRouter arma el router de la API usando el enrutador nativo de net/http
// (Go 1.22+), con los patrones por método y los path params (id).
func NewRouter(service application.ProductService) http.Handler {
	mux := http.NewServeMux()
	products := handler.NewProductHandler(service)

	// Healthcheck para readiness/liveness probes.
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Documentación interactiva (Swagger UI) y spec OpenAPI embebido.
	mux.HandleFunc("GET /docs", docs.SwaggerUI)
	mux.HandleFunc("GET /openapi.yaml", docs.OpenAPISpec)

	mux.HandleFunc("GET /products", products.GetAll)
	mux.HandleFunc("POST /products", products.Create)
	mux.HandleFunc("GET /products/{id}", products.GetByID)
	mux.HandleFunc("PUT /products/{id}", products.Update)
	mux.HandleFunc("DELETE /products/{id}", products.Delete)

	// Los middlewares se aplican en cadena: Recover envuelve a Logging, que
	// envuelve al router. Recover queda más afuera para capturar cualquier panic.
	return middleware.Recover(middleware.Logging(mux))
}
