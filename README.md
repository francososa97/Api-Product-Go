# Product API

API REST de productos escrita en **Go**, diseñada como muestrario de **Clean Architecture** aplicada de forma estricta: dependencias que apuntan siempre hacia el dominio, lógica de negocio aislada del transporte y de la base de datos, e inyección de dependencias en cada capa.

Corre **sin dependencias externas** gracias a un repositorio in-memory por defecto, y puede cambiarse a **MongoDB** con una sola variable de entorno.

---

## ✨ Características

- **Clean Architecture** estricta (domain → application → infrastructure → presentation).
- **Documentación interactiva** con Swagger UI (OpenAPI 3) embebida en el binario.
- **CRUD completo** de productos con ordenamiento por precio.
- **Dos motores de persistencia** intercambiables: in-memory (default) y MongoDB.
- **Router nativo** de `net/http` (Go 1.22+), sin frameworks de routing.
- **Una sola dependencia externa**: el driver oficial de MongoDB.
- Respuestas de **error consistentes** (`{ error, message, statusCode }`).
- **Graceful shutdown**, timeouts de request y contextos con cancelación.
- **Middlewares** de logging y recuperación de panics.
- **Validación** de dominio y generación de IDs (UUID v4) sin librerías de terceros.
- **Tests unitarios** con mocks en cada capa.
- Listo para **Docker** y **docker-compose**.

---

## 🏗️ Arquitectura

El proyecto sigue Clean Architecture: las flechas de dependencia apuntan siempre hacia adentro. El **dominio no conoce** ni HTTP ni MongoDB.

```text
┌─────────────────────────────────────────────────────────────┐
│  presentation (HTTP)   handlers · router · middleware         │
│        │  depende de                                          │
│        ▼                                                       │
│  application           casos de uso (ProductService)          │
│        │  depende de                                          │
│        ▼                                                       │
│  domain                Product · reglas · puerto Repository   │  ← núcleo, sin dependencias
│        ▲  implementa                                          │
│        │                                                       │
│  infrastructure        repos (memory · mongo) · config        │
└─────────────────────────────────────────────────────────────┘
```

### Estructura de carpetas

```text
.
├── cmd/
│   └── api/
│       └── main.go                 # Composition root + graceful shutdown
├── internal/
│   ├── domain/                     # Entidades, reglas de negocio, puertos
│   │   ├── product.go              # Entidad Product + validación
│   │   ├── id.go                   # Generación de UUID v4 (crypto/rand)
│   │   ├── errors.go               # Errores de dominio
│   │   └── repository.go           # Puerto ProductRepository (interfaz)
│   ├── application/                # Casos de uso
│   │   └── product_service.go
│   ├── infrastructure/             # Implementaciones concretas
│   │   ├── config/                 # Carga y validación de env vars
│   │   └── persistence/
│   │       ├── memory/             # Repositorio in-memory (default)
│   │       └── mongo/              # Repositorio MongoDB
│   └── presentation/
│       └── http/                   # Handlers, router, middleware, respuestas
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

> Todo el código vive bajo `internal/`, por lo que no puede ser importado desde otros módulos: la API expone una interfaz HTTP, no una librería.

---

## 🚀 Inicio rápido

Requisitos: **Go 1.22+** (y opcionalmente Docker).

```bash
git clone https://github.com/francososa97/product-api.git
cd product-api
go mod tidy   # resuelve y verifica las dependencias
go run ./cmd/api
```

La API queda escuchando en `http://localhost:8080` con el repositorio **in-memory**. No necesitás base de datos para probarla.

### Con Makefile

```bash
make run     # ejecuta la API
make test    # corre los tests con race detector
make cover   # genera reporte de cobertura HTML
make help    # lista todos los comandos
```

### Con Docker + MongoDB

```bash
docker compose up --build
```

Levanta la API (con `DB_DRIVER=mongo`) junto a una instancia de MongoDB, lista en `http://localhost:8080`.

---

## ⚙️ Configuración

Todas las opciones se leen de variables de entorno (ver [`.env.example`](.env.example)):

| Variable           | Default                     | Descripción                                      |
|--------------------|-----------------------------|--------------------------------------------------|
| `APP_PORT`         | `8080`                      | Puerto HTTP.                                      |
| `DB_DRIVER`        | `memory`                    | Motor de persistencia: `memory` o `mongo`.       |
| `MONGO_URI`        | `mongodb://localhost:27017` | URI de conexión (solo con `DB_DRIVER=mongo`).    |
| `MONGO_DB`         | `productsdb`                | Nombre de la base de datos.                      |
| `MONGO_COLLECTION` | `products`                  | Nombre de la colección.                          |
| `REQUEST_TIMEOUT`  | `5s`                        | Timeout de lectura/escritura de cada request.    |
| `SHUTDOWN_TIMEOUT` | `10s`                       | Tiempo máximo para el apagado ordenado.          |

Para usar MongoDB localmente sin Docker:

```bash
DB_DRIVER=mongo MONGO_URI=mongodb://localhost:27017 go run ./cmd/api
```

---

## 📖 Documentación interactiva

Con la API corriendo, abrí **Swagger UI** en el navegador:

```text
http://localhost:8080/docs
```

Desde ahí podés explorar todos los endpoints y probarlos con el botón **Try it out**. El spec OpenAPI 3 crudo se sirve en:

```text
http://localhost:8080/openapi.yaml
```

> El spec está embebido en el binario con `go:embed` (no hace falta servir archivos aparte) y la fuente vive en [`internal/presentation/http/docs/openapi.yaml`](internal/presentation/http/docs/openapi.yaml).

---

## 📡 Endpoints

Base URL: `http://localhost:8080`

| Método   | Ruta              | Descripción                                   | Body | Respuesta |
|----------|-------------------|-----------------------------------------------|------|-----------|
| `GET`    | `/health`         | Healthcheck                                   | —    | `200`     |
| `GET`    | `/docs`           | Documentación interactiva (Swagger UI)        | —    | `200`     |
| `GET`    | `/openapi.yaml`   | Spec OpenAPI 3                                 | —    | `200`     |
| `GET`    | `/products`       | Lista productos (`?sortByPriceAsc=true`)      | —    | `200`     |
| `GET`    | `/products/{id}`  | Obtiene un producto por ID                    | —    | `200` / `404` |
| `POST`   | `/products`       | Crea un producto                              | sí   | `201` / `400` |
| `PUT`    | `/products/{id}`  | Actualiza un producto                         | sí   | `200` / `400` / `404` |
| `DELETE` | `/products/{id}`  | Elimina un producto                           | —    | `204` / `404` |

### Modelo

```json
{
  "id": "8f14e45f-ceea-467d-9a1b-2c3d4e5f6a7b",
  "name": "Teclado mecánico",
  "price": 149.99
}
```

El `id` lo genera el servidor (UUID v4); el cliente solo envía `name` y `price`.

### Formato de error

Todas las respuestas de error comparten el mismo formato:

```json
{
  "error": "Not Found",
  "message": "producto no encontrado",
  "statusCode": 404
}
```

---

## 🧪 Ejemplos con cURL

```bash
# Crear un producto
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Teclado mecánico","price":149.99}'

# Listar productos ordenados por precio ascendente
curl "http://localhost:8080/products?sortByPriceAsc=true"

# Obtener un producto por ID
curl http://localhost:8080/products/<id>

# Actualizar un producto
curl -X PUT http://localhost:8080/products/<id> \
  -H "Content-Type: application/json" \
  -d '{"name":"Teclado mecánico RGB","price":179.99}'

# Eliminar un producto
curl -X DELETE http://localhost:8080/products/<id>
```

---

## 🧱 Decisiones de diseño

- **El dominio no depende de nada.** `Product`, sus reglas (`Validate`, `Update`) y el puerto `ProductRepository` viven en `domain` y no importan ni HTTP ni Mongo. Las implementaciones concretas dependen del dominio, nunca al revés.
- **Inyección de dependencias.** El *composition root* (`cmd/api/main.go`) decide qué implementación usar y la inyecta; cada capa recibe interfaces, no concreciones.
- **Sin framework de routing.** Se usa el `ServeMux` de Go 1.22, que ya soporta patrones por método (`GET /products/{id}`) y path params. Menos dependencias, menos superficie de mantenimiento.
- **IDs sin dependencias.** Los UUID v4 se generan con `crypto/rand` de la biblioteca estándar.
- **Errores de dominio tipados.** El dominio devuelve errores semánticos (`ErrProductNotFound`, etc.) y la capa HTTP los traduce a códigos de estado con `errors.Is`, sin filtrar detalles internos al cliente.

---

## ✅ Tests

```bash
go test ./...           # todos los tests
go test -race ./...     # con detección de race conditions
make cover              # reporte de cobertura HTML
```

Hay tests unitarios para la lógica de dominio, los casos de uso (con un mock del repositorio) y el repositorio in-memory.

---

## 🛠️ Stack

- **Go 1.22+**
- **net/http** (router nativo)
- **MongoDB** vía [`mongo-driver`](https://github.com/mongodb/mongo-go-driver) (opcional)
- **Docker** / **docker-compose**

---

## 📄 Licencia

[MIT](LICENSE) © Franco Sosa
