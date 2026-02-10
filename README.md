# Clothes Shop API

A RESTful API for a clothes shop built with Golang, Gin, and PostgreSQL.

## Features

- User management
- Product catalog with categories
- Shopping cart
- Order management
- Database migrations
- Docker support

## Prerequisites

- Go 1.21+
- Docker and Docker Compose (on Windows, ensure Docker Desktop is running with administrative privileges)
- golang-migrate CLI (installed via `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`)

## Setup

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd clothes-shop-api
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Start PostgreSQL with Docker:
   ```bash
   docker-compose up -d
   ```

4. Run migrations:
   ```bash
   migrate -database "postgres://postgres:123456789@localhost:5432/ClothesShopDB?sslmode=disable" -path migrations up
   ```

   **Note:** If you get "database does not exist" error, wait a moment for Docker to fully initialize the database, or check that Docker is running properly.

5. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

The API will be available at `http://localhost:8080`.

## Troubleshooting

### Windows Docker Issues
If you encounter errors like "database does not exist" or Docker connection issues on Windows:

1. **Ensure Docker Desktop is running with administrative privileges:**
   - Right-click Docker Desktop and select "Run as administrator"
   - Or start PowerShell/Command Prompt as administrator and run Docker commands

2. **Check Docker status:**
   ```bash
   docker --version
   docker ps
   ```

3. **If Docker is not available, you can set up PostgreSQL locally:**
   - Install PostgreSQL from https://www.postgresql.org/download/windows/
   - Create database named "ClothesShopDB"
   - Update the connection string in migration commands to match your local setup

4. **Verify database creation:**
   ```bash
   docker exec -it clothes-shop-api-postgres-1 psql -U postgres -d ClothesShopDB -c "\l"
   ```

## API Endpoints

- `GET /health` - Health check

## Environment Variables

Create a `.env` file in the root directory:

```
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ClothesShopDB
DB_USER=postgres
DB_PASSWORD=123456789
```

## Project Structure

```
clothes-shop-api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── db/
│   │   └── postgres.go
│   ├── models/
│   │   ├── user.go
│   │   ├── category.go
│   │   ├── product.go
│   │   ├── cart.go
│   │   └── order.go
│   ├── handlers/
│   ├── repositories/
│   └── routes/
├── migrations/
│   ├── 000001_init_tables.up.sql
│   └── 000001_init_tables.down.sql
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```

## Migration

To create a new migration:
```bash
migrate create -ext sql -dir migrations -seq <migration_name>
```

To rollback:
```bash
migrate -database "postgres://postgres:123456789@localhost:5432/mydb?sslmode=disable" -path migrations down

Sinh ra tai lieu swagger:
```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go