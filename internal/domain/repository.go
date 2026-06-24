package domain

import "context"

// ProductRepository es el puerto de persistencia: lo define el dominio y lo
// implementan las capas externas (in-memory, MongoDB, etc.). Así la lógica de
// negocio depende de una abstracción y nunca de una tecnología concreta.
//
// Implementaciones esperadas:
//   - GetByID debe devolver ErrProductNotFound si el producto no existe.
//   - Todos los métodos reciben context para soportar timeouts y cancelación.
type ProductRepository interface {
	GetAll(ctx context.Context, sortByPriceAsc bool) ([]Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
}
