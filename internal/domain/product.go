package domain

import "strings"

// Product es la entidad central del dominio. No conoce la base de datos ni el
// transporte HTTP: solo modela un producto y sus reglas de negocio.
type Product struct {
	ID    string  `json:"id" bson:"_id"`
	Name  string  `json:"name" bson:"name"`
	Price float64 `json:"price" bson:"price"`
}

// NewProduct construye un producto válido y le asigna una identidad propia.
// Centralizar la creación acá garantiza que nunca exista un Product inválido
// o sin ID dando vueltas por las capas superiores.
func NewProduct(name string, price float64) (*Product, error) {
	id, err := newID()
	if err != nil {
		return nil, err
	}

	p := &Product{
		ID:    id,
		Name:  strings.TrimSpace(name),
		Price: price,
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}

// Validate aplica las invariantes del producto. Es la única fuente de verdad
// sobre qué hace válido a un producto, reutilizada en la creación y la edición.
func (p *Product) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return ErrEmptyName
	}
	if p.Price < 0 {
		return ErrNegativePrice
	}
	return nil
}

// Update aplica cambios sobre un producto existente conservando su identidad
// y revalidando el resultado antes de persistirlo.
func (p *Product) Update(name string, price float64) error {
	p.Name = strings.TrimSpace(name)
	p.Price = price
	return p.Validate()
}
