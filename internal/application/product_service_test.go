package application

import (
	"context"
	"errors"
	"testing"

	"github.com/francososa97/product-api/internal/domain"
)

// mockRepository es un doble de prueba del puerto domain.ProductRepository.
// Permite testear los casos de uso de forma aislada, sin base de datos real.
type mockRepository struct {
	getByIDFunc func(ctx context.Context, id string) (*domain.Product, error)
	createFunc  func(ctx context.Context, p *domain.Product) error
	updateFunc  func(ctx context.Context, p *domain.Product) error
	deleteFunc  func(ctx context.Context, id string) error
}

func (m *mockRepository) GetAll(context.Context, bool) ([]domain.Product, error) { return nil, nil }
func (m *mockRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	return m.getByIDFunc(ctx, id)
}
func (m *mockRepository) Create(ctx context.Context, p *domain.Product) error {
	return m.createFunc(ctx, p)
}
func (m *mockRepository) Update(ctx context.Context, p *domain.Product) error {
	return m.updateFunc(ctx, p)
}
func (m *mockRepository) Delete(ctx context.Context, id string) error { return m.deleteFunc(ctx, id) }

func TestCreateProduct_shouldPersistGeneratedProductWhenInputIsValid(t *testing.T) {
	var saved *domain.Product
	repo := &mockRepository{
		createFunc: func(_ context.Context, p *domain.Product) error {
			saved = p
			return nil
		},
	}
	svc := NewProductService(repo)

	got, err := svc.CreateProduct(context.Background(), CreateProductInput{Name: "Monitor", Price: 200})
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo: %v", err)
	}
	if got.ID == "" {
		t.Error("se esperaba un ID generado por el servicio")
	}
	if saved == nil || saved.ID != got.ID {
		t.Error("el producto creado debería haberse persistido")
	}
}

func TestCreateProduct_shouldFailWhenInputIsInvalid(t *testing.T) {
	repo := &mockRepository{
		createFunc: func(context.Context, *domain.Product) error {
			t.Fatal("no debería persistirse un producto inválido")
			return nil
		},
	}
	svc := NewProductService(repo)

	if _, err := svc.CreateProduct(context.Background(), CreateProductInput{Name: "", Price: 10}); !errors.Is(err, domain.ErrEmptyName) {
		t.Errorf("se esperaba ErrEmptyName, se obtuvo: %v", err)
	}
}

func TestUpdateProduct_shouldReturnNotFoundWhenProductDoesNotExist(t *testing.T) {
	repo := &mockRepository{
		getByIDFunc: func(context.Context, string) (*domain.Product, error) {
			return nil, domain.ErrProductNotFound
		},
	}
	svc := NewProductService(repo)

	if _, err := svc.UpdateProduct(context.Background(), "x", UpdateProductInput{Name: "Nuevo", Price: 1}); !errors.Is(err, domain.ErrProductNotFound) {
		t.Errorf("se esperaba ErrProductNotFound, se obtuvo: %v", err)
	}
}

func TestUpdateProduct_shouldPreserveIDWhenProductExists(t *testing.T) {
	existing, _ := domain.NewProduct("Viejo", 10)
	var updated *domain.Product
	repo := &mockRepository{
		getByIDFunc: func(context.Context, string) (*domain.Product, error) { return existing, nil },
		updateFunc: func(_ context.Context, p *domain.Product) error {
			updated = p
			return nil
		},
	}
	svc := NewProductService(repo)

	got, err := svc.UpdateProduct(context.Background(), existing.ID, UpdateProductInput{Name: "Nuevo", Price: 30})
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo: %v", err)
	}
	if got.ID != existing.ID || updated.ID != existing.ID {
		t.Error("el ID debería preservarse al actualizar")
	}
	if got.Name != "Nuevo" || got.Price != 30 {
		t.Errorf("los cambios no se aplicaron: %+v", got)
	}
}

func TestDeleteProduct_shouldReturnNotFoundWhenProductDoesNotExist(t *testing.T) {
	repo := &mockRepository{
		getByIDFunc: func(context.Context, string) (*domain.Product, error) {
			return nil, domain.ErrProductNotFound
		},
	}
	svc := NewProductService(repo)

	if err := svc.DeleteProduct(context.Background(), "x"); !errors.Is(err, domain.ErrProductNotFound) {
		t.Errorf("se esperaba ErrProductNotFound, se obtuvo: %v", err)
	}
}
