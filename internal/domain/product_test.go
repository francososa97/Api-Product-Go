package domain

import (
	"errors"
	"testing"
)

func TestNewProduct_shouldCreateValidProductWhenInputIsValid(t *testing.T) {
	p, err := NewProduct("  Teclado  ", 99.9)
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo: %v", err)
	}
	if p.ID == "" {
		t.Error("se esperaba un ID generado, se obtuvo vacío")
	}
	if p.Name != "Teclado" {
		t.Errorf("se esperaba el nombre trimmeado %q, se obtuvo %q", "Teclado", p.Name)
	}
}

func TestNewProduct_shouldFailWhenNameIsEmpty(t *testing.T) {
	if _, err := NewProduct("   ", 10); !errors.Is(err, ErrEmptyName) {
		t.Errorf("se esperaba ErrEmptyName, se obtuvo: %v", err)
	}
}

func TestNewProduct_shouldFailWhenPriceIsNegative(t *testing.T) {
	if _, err := NewProduct("Mouse", -1); !errors.Is(err, ErrNegativePrice) {
		t.Errorf("se esperaba ErrNegativePrice, se obtuvo: %v", err)
	}
}

func TestUpdate_shouldKeepIDAndApplyChangesWhenValid(t *testing.T) {
	p, _ := NewProduct("Mouse", 10)
	originalID := p.ID

	if err := p.Update("Mouse Pro", 25); err != nil {
		t.Fatalf("no se esperaba error, se obtuvo: %v", err)
	}
	if p.ID != originalID {
		t.Error("el ID no debería cambiar al actualizar")
	}
	if p.Name != "Mouse Pro" || p.Price != 25 {
		t.Errorf("los cambios no se aplicaron: %+v", p)
	}
}

func TestUpdate_shouldFailWhenResultIsInvalid(t *testing.T) {
	p, _ := NewProduct("Mouse", 10)
	if err := p.Update("Mouse", -5); !errors.Is(err, ErrNegativePrice) {
		t.Errorf("se esperaba ErrNegativePrice, se obtuvo: %v", err)
	}
}

func TestNewID_shouldGenerateUniqueIDs(t *testing.T) {
	seen := make(map[string]struct{})
	for i := 0; i < 1000; i++ {
		id, err := newID()
		if err != nil {
			t.Fatalf("no se esperaba error generando ID: %v", err)
		}
		if _, dup := seen[id]; dup {
			t.Fatalf("ID duplicado generado: %s", id)
		}
		seen[id] = struct{}{}
	}
}
