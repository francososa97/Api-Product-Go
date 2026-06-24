// Package memory provee una implementación en memoria de domain.ProductRepository.
// Es el motor por defecto: permite levantar la API sin dependencias externas,
// ideal para desarrollo, demos y tests.
package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/francososa97/product-api/internal/domain"
)

// ProductRepository almacena los productos en un mapa protegido por un mutex,
// por lo que es seguro para uso concurrente.
type ProductRepository struct {
	mu       sync.RWMutex
	products map[string]domain.Product
}

// NewProductRepository crea un repositorio en memoria vacío.
func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		products: make(map[string]domain.Product),
	}
}

func (r *ProductRepository) GetAll(_ context.Context, sortByPriceAsc bool) ([]domain.Product, error) {
	r.mu.RLock()
	products := make([]domain.Product, 0, len(r.products))
	for _, p := range r.products {
		products = append(products, p)
	}
	r.mu.RUnlock()

	sort.Slice(products, func(i, j int) bool {
		if sortByPriceAsc {
			return products[i].Price < products[j].Price
		}
		return products[i].Price > products[j].Price
	})

	return products, nil
}

func (r *ProductRepository) GetByID(_ context.Context, id string) (*domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, ok := r.products[id]
	if !ok {
		return nil, domain.ErrProductNotFound
	}
	return &product, nil
}

func (r *ProductRepository) Create(_ context.Context, product *domain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.products[product.ID] = *product
	return nil
}

func (r *ProductRepository) Update(_ context.Context, product *domain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[product.ID]; !ok {
		return domain.ErrProductNotFound
	}
	r.products[product.ID] = *product
	return nil
}

func (r *ProductRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[id]; !ok {
		return domain.ErrProductNotFound
	}
	delete(r.products, id)
	return nil
}
