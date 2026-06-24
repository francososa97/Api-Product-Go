.PHONY: help run build test cover lint tidy docker-up docker-down clean

BINARY := product-api
PKG := ./cmd/api

help: ## Muestra esta ayuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-14s\033[0m %s\n", $$1, $$2}'

run: ## Ejecuta la API (repositorio in-memory por defecto)
	go run $(PKG)

build: ## Compila el binario en ./bin
	go build -o bin/$(BINARY) $(PKG)

test: ## Corre los tests con race detector
	go test -race ./...

cover: ## Corre los tests y genera reporte de cobertura HTML
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Ejecuta go vet
	go vet ./...

tidy: ## Ordena y verifica las dependencias
	go mod tidy

docker-up: ## Levanta API + MongoDB con docker compose
	docker compose up --build

docker-down: ## Detiene y elimina los contenedores
	docker compose down

clean: ## Elimina binarios y artefactos de cobertura
	rm -rf bin coverage.out coverage.html
