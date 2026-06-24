// Package handler contiene los handlers HTTP de productos. Su única
// responsabilidad es traducir entre HTTP y los casos de uso: decodificar la
// petición, invocar el servicio y serializar la respuesta.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/francososa97/product-api/internal/application"
	"github.com/francososa97/product-api/internal/presentation/http/response"
)

// ProductHandler agrupa los handlers de productos sobre un caso de uso.
type ProductHandler struct {
	service application.ProductService
}

// NewProductHandler crea el handler inyectando el servicio (interfaz).
func NewProductHandler(service application.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// productRequest es el cuerpo aceptado para crear y actualizar productos.
type productRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// GetAll responde la lista de productos, opcionalmente ordenada por precio
// ascendente con ?sortByPriceAsc=true.
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	sortByPriceAsc := r.URL.Query().Get("sortByPriceAsc") == "true"

	products, err := h.service.GetAllProducts(r.Context(), sortByPriceAsc)
	if err != nil {
		response.FromDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, products)
}

// GetByID responde un producto por su ID.
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	product, err := h.service.GetProductByID(r.Context(), r.PathValue("id"))
	if err != nil {
		response.FromDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, product)
}

// Create da de alta un producto y devuelve el recurso creado con su ID.
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	req, ok := decode(w, r)
	if !ok {
		return
	}

	product, err := h.service.CreateProduct(r.Context(), application.CreateProductInput{
		Name:  req.Name,
		Price: req.Price,
	})
	if err != nil {
		response.FromDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, product)
}

// Update modifica un producto existente y devuelve el recurso actualizado.
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	req, ok := decode(w, r)
	if !ok {
		return
	}

	product, err := h.service.UpdateProduct(r.Context(), r.PathValue("id"), application.UpdateProductInput{
		Name:  req.Name,
		Price: req.Price,
	})
	if err != nil {
		response.FromDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, product)
}

// Delete elimina un producto por su ID.
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteProduct(r.Context(), r.PathValue("id")); err != nil {
		response.FromDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// decode lee el cuerpo JSON de la petición en un productRequest. Si el JSON es
// inválido responde 400 y devuelve ok=false para que el handler corte.
func decode(w http.ResponseWriter, r *http.Request) (productRequest, bool) {
	var req productRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "cuerpo JSON inválido: "+err.Error())
		return productRequest{}, false
	}
	return req, true
}
