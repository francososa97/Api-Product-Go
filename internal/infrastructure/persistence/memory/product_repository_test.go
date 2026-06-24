package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/francososa97/product-api/internal/domain"
)

func seed(t *testing.T, repo *ProductRepository, name string, price float64) *domain.Product {
	t.Helper()
	p, err := domain.NewProduct(name, price)
	if err != nil {
		t.Fatalf("no se pudo crear el producto de prueba: %v", err)
	}
	if err := repo.Create(context.Background(), p); err != nil {
		t.Fatalf("no se pudo persistir el producto de prueba: %v", err)
	}
	return p
}

func TestGetByID_shouldReturnNotFoundWhenAbsent(t *testing.T) {
	repo := NewProductRepository()
	if _, err := repo.GetByID(context.Background(), "inexistente"); !errors.Is(err, domain.ErrProductNotFound) {
		t.Errorf("se esperaba ErrProductNotFound, se obtuvo: %v", err)
	}
}

func TestCreateAndGetByID_shouldRoundTrip(t *testing.T) {
	repo := NewProductRepository()
	created := seed(t, repo, "Webcam", 50)

	got, err := repo.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo: %v", err)
	}
	if got.ID != created.ID || got.Name != "Webcam" {
		t.Errorf("se esperaba %+v, se obtuvo %+v", created, got)
	}
}

func TestGetAll_shouldSortByPrice(t *testing.T) {
	repo := NewProductRepository()
	seed(t, repo, "Caro", 300)
	seed(t, repo, "Barato", 100)
	seed(t, repo, "Medio", 200)

	asc, _ := repo.GetAll(context.Background(), true)
	if asc[0].Price != 100 || asc[2].Price != 300 {
		t.Errorf("orden ascendente incorrecto: %+v", asc)
	}

	desc, _ := repo.GetAll(context.Background(), false)
	if desc[0].Price != 300 || desc[2].Price != 100 {
		t.Errorf("orden descendente incorrecto: %+v", desc)
	}
}

func TestUpdate_shouldFailWhenAbsent(t *testing.T) {
	repo := NewProductRepository()
	ghost, _ := domain.NewProduct("Fantasma", 1)
	if err := repo.Update(context.Background(), ghost); !errors.Is(err, domain.ErrProductNotFound) {
		t.Errorf("se esperaba ErrProductNotFound, se obtuvo: %v", err)
	}
}

func TestDelete_shouldRemoveExistingProduct(t *testing.T) {
	repo := NewProductRepository()
	created := seed(t, repo, "Borrable", 10)

	if err := repo.Delete(context.Background(), created.ID); err != nil {
		t.Fatalf("no se esperaba error al borrar: %v", err)
	}
	if _, err := repo.GetByID(context.Background(), created.ID); !errors.Is(err, domain.ErrProductNotFound) {
		t.Errorf("el producto debería haberse borrado, error: %v", err)
	}
}

func TestDelete_shouldFailWhenAbsent(t *testing.T) {
	repo := NewProductRepository()
	if err := repo.Delete(context.Background(), "inexistente"); !errors.Is(err, domain.ErrProductNotFound) {
		t.Errorf("se esperaba ErrProductNotFound, se obtuvo: %v", err)
	}
}
