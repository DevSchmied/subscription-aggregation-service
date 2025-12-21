# Subscription Aggregation Service

A REST service for managing and aggregating users’ online subscriptions.

The project demonstrates clean architecture, input validation, error handling, PostgreSQL integration, Swagger API documentation, Docker-based deployment, and a basic CI pipeline.

--- 

## Task Description

Design and implement a REST service for aggregating data about users’ online subscriptions.

---

## Functional Requirements

- CRUDL operations for subscription records  
- HTTP endpoint for calculating the total subscription cost for a selected period  
- Filtering by:
  - `user_id`
  - `service_name`
- PostgreSQL used as the database  
- Database migrations for schema initialization  
- Application logging  
- Configuration via `.env` file  
- Swagger API documentation  
- Service startup using `docker compose`

---

## Notes

- User existence validation is **out of scope**
- Subscription price is an **integer value**
- Date format: **MM-YYYY**

---

## Domain Model
### Subscription

- `id` — UUID  
- `service_name` — subscription service name  
- `price` — monthly price (integer)  
- `user_id` — UUID  
- `start_date` — MM-YYYY  
- `end_date` — optional, MM-YYYY  
- `created_at` — creation timestamp  
- `updated_at` — last update timestamp  

---

## Technology Stack

- Go 1.24  
- Gin — HTTP framework  
- PostgreSQL 16  
- pgx / pgxpool — database access  
- Swagger (swaggo) — API documentation  
- Docker & Docker Compose  
- GitHub Actions — CI  
- godotenv — configuration loading  

---

## Project Structure

- `cmd/app` — application entry point  
- `cmd/smoke` — smoke test runner  
- `internal/config` — configuration loading and validation  
- `internal/domain` — domain entities  
- `internal/http/handlers` — HTTP handlers and unit tests  
- `internal/http/router` — Gin router configuration  
- `internal/storage/postgres` — PostgreSQL repository  
- `internal/utils` — date handling utilities and unit tests
- `migrations` — SQL migrations  
- `docs` — generated Swagger documentation  
- `docker-compose.yml` — Docker Compose configuration  
- `Dockerfile` — application Docker image  
- `.env.example` — environment configuration example  
- `README.md` — project documentation  

---

## Running the Service
### 1 Environment Setup

Create a `.env` file based on the example:

```bash
cp .env.example .env
```

### Example `.env` Configuration

```bash
APP_PORT=8080

DB_PORT=5432
DB_NAME=subscriptions
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSLMODE=disable
```

## 2 Run with Docker Compose

Start the application using Docker Compose:

```bash
docker compose up --build
```

**During Startup** 
- PostgreSQL database is started
- Database migrations are applied automatically
- The application becomes available on port 8080

---

## Swagger Documentation

http://localhost:8080/swagger/index.html

The documentation is fully generated from annotations in the code.

---

## API Endpoints
### Subscriptions

- `POST /api/subscriptions` — create a subscription  
- `GET /api/subscriptions/{id}` — get subscription by ID  
- `PUT /api/subscriptions/{id}` — update subscription  
- `DELETE /api/subscriptions/{id}` — delete subscription  
- `GET /api/subscriptions` — list subscriptions  

### Aggregation

- `GET /api/subscriptions/total` — calculate total subscription cost

#### Aggregation Parameters

- `start_date` (required) — MM-YYYY  
- `end_date` (required) — MM-YYYY  
- `user_id` (optional)  
- `service_name` (optional)  

---

## Testing

Run all tests:

```bash
go test ./...
```

The project is covered by unit tests, including:
- date parsing and formatting
- input validation
- error-to-HTTP mapping
- aggregation helper functions
- HTTP response formatting

---

## CI (Continuous Integration)

GitHub Actions runs on every push and pull request to the `main` branch.

Executed steps:
- `go test ./...`  
- `go vet ./...`  
- `go build ./cmd/app`  

CI configuration is located in:

.github/workflows/ci.yml

---

## Logging and Error Handling

The application logs:

- validation errors  
- database errors  
- successful CRUD operations  

Centralized error mapping is used:

- `404` — not found  
- `504` — timeout  
- `500` — database error  
