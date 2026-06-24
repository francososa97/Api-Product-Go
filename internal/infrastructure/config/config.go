// Package config carga y valida la configuración de la aplicación desde
// variables de entorno al iniciar, fallando rápido si algo es inválido.
package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Driver identifica el motor de persistencia a utilizar.
type Driver string

const (
	DriverMemory Driver = "memory"
	DriverMongo  Driver = "mongo"
)

// Config agrupa toda la configuración de la API.
type Config struct {
	Port            string
	Driver          Driver
	MongoURI        string
	MongoDB         string
	MongoCollection string
	RequestTimeout  time.Duration
	ShutdownTimeout time.Duration
}

// Load lee las variables de entorno, aplica valores por defecto razonables y
// valida el resultado. Devuelve error en lugar de hacer panic para que el
// llamador decida cómo reportar el fallo de arranque.
func Load() (*Config, error) {
	cfg := &Config{
		Port:            getEnv("APP_PORT", "8080"),
		Driver:          Driver(strings.ToLower(getEnv("DB_DRIVER", string(DriverMemory)))),
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:         getEnv("MONGO_DB", "productsdb"),
		MongoCollection: getEnv("MONGO_COLLECTION", "products"),
		RequestTimeout:  getEnvDuration("REQUEST_TIMEOUT", 5*time.Second),
		ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	switch c.Driver {
	case DriverMemory, DriverMongo:
		// válidos
	default:
		return fmt.Errorf("DB_DRIVER inválido: %q (use %q o %q)", c.Driver, DriverMemory, DriverMongo)
	}

	if c.Driver == DriverMongo {
		if c.MongoURI == "" {
			return fmt.Errorf("MONGO_URI es obligatorio cuando DB_DRIVER=%s", DriverMongo)
		}
		if c.MongoDB == "" {
			return fmt.Errorf("MONGO_DB es obligatorio cuando DB_DRIVER=%s", DriverMongo)
		}
	}

	return nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return v
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
