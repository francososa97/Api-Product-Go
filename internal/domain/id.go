package domain

import (
	"crypto/rand"
	"fmt"
)

// newID genera un identificador único con formato UUID v4 (RFC 4122).
// Se implementa sobre crypto/rand (biblioteca estándar) para no acoplar el
// dominio a ninguna dependencia externa de generación de IDs.
func newID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("no se pudo generar el ID del producto: %w", err)
	}

	b[6] = (b[6] & 0x0f) | 0x40 // versión 4
	b[8] = (b[8] & 0x3f) | 0x80 // variante RFC 4122

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
