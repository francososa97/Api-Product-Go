package domain

import "errors"

// Errores de dominio. Las capas externas deben compararlos con errors.Is
// y traducirlos a la respuesta que corresponda (por ejemplo, un código HTTP),
// sin que el dominio conozca nada del transporte.
var (
	// ErrProductNotFound se devuelve cuando no existe un producto con el ID pedido.
	ErrProductNotFound = errors.New("producto no encontrado")

	// ErrEmptyName se devuelve cuando el nombre del producto está vacío.
	ErrEmptyName = errors.New("el nombre del producto no puede estar vacío")

	// ErrNegativePrice se devuelve cuando el precio del producto es negativo.
	ErrNegativePrice = errors.New("el precio del producto no puede ser negativo")
)
