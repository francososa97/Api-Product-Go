// Package application contiene los casos de uso de la API. Orquesta el dominio
// y la persistencia, pero no sabe nada de HTTP ni de la base de datos concreta:
// solo depende del puerto domain.ProductRepository.
package application

import (
	"context"

	"github.com/francososa97/product-api/internal/domain"
)

// CreateProductInput son los datos necesarios para dar de alta un producto.
type CreateProductInput struct {
	Name  string
	Price float64
}

// UpdateProductInput son los datos modificables de un producto existente.
type UpdateProductInput struct {
	Name  string
	Price float64
}

// ProductService define los casos de uso disponibles sobre productos. La capa
// de presentación depende de esta interfaz, no de la implementación concreta.
type ProductService interface {
	GetAllProducts(ctx context.Context, sortByPriceAsc bool) ([]domain.Product, error)
	GetProductByID(ctx context.Context, id string) (*domain.Product, error)
	CreateProduct(ctx context.Context, in CreateProductInput) (*domain.Product, error)
	UpdateProduct(ctx context.Context, id string, in UpdateProductInput) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id string) error
}

// productService es la implementación por defecto de los casos de uso.
type productService struct {
	repo domain.ProductRepository
}

// NewProductService crea el servicio inyectando el repositorio (inversión de
// dependencias: recibe la abstracción, no una implementación concreta).
func NewProductService(repo domain.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAllProducts(ctx context.Context, sortByPriceAsc bool) ([]domain.Product, error) {
	return s.repo.GetAll(ctx, sortByPriceAsc)
}

func (s *productService) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

// CreateProduct construye la entidad (validando e identificando) y la persiste.
func (s *productService) CreateProduct(ctx context.Context, in CreateProductInput) (*domain.Product, error) {
	product, err := domain.NewProduct(in.Name, in.Price)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

// UpdateProduct recupera el producto existente para conservar su identidad,
// le aplica los cambios validados y lo persiste. Si no existe, propaga
// domain.ErrProductNotFound.
func (s *productService) UpdateProduct(ctx context.Context, id string, in UpdateProductInput) (*domain.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := product.Update(in.Name, in.Price); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// DeleteProduct verifica la existencia antes de borrar para devolver un
// resultado coherente (404 en lugar de un borrado silencioso de nada).
func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
