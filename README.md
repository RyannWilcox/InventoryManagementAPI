# Inventory Management API

_Udacity project #2_

A RESTful backend server for managing an inventory of items, built in Go with [Gin](https://github.com/gin-gonic/gin) and [GORM](https://gorm.io/). Users can perform full CRUD operations on inventory items each identified by a UUID. The API supports pagination, sorting, filtering, and rate limiting.

## Features

- **CRUD operations** — list, retrieve, create, update, and delete inventory items
- **UUID primary keys** — every item is assigned a unique, non-sequential identifier on creation
- **Pagination** — retrieve large inventories in manageable chunks via `limit`/`offset` query params
- **Sorting** — sort results by `name`, `stock`, or `price`, in ascending or descending order
- **Filtering** — filter items by name or minimum stock level
- **Rate limiting** — a token-bucket limiter protects the API from excessive request volume
- **Centralized error handling** — consistent JSON error responses across all endpoints
- **Swagger/OpenAPI docs** — interactive API documentation generated from code comments
- **Database seeding** — automatically populates sample data on first run

## Tech Stack

| Component        | Technology                                              |
|-------------------|----------------------------------------------------------|
| Language          | Go 1.25                                                  |
| Web framework     | [Gin](https://github.com/gin-gonic/gin)                  |
| ORM               | [GORM](https://gorm.io/) with the PostgreSQL driver       |
| Database          | PostgreSQL                                                |
| API documentation | [swaggo/swag](https://github.com/swaggo/swag) + gin-swagger |
| Config            | [godotenv](https://github.com/joho/godotenv) (`.env` files) |

## Prerequisites

Before you begin, make sure you have the following installed:

- **[Go](https://go.dev/dl/)** 1.25 or later — verify with:
  ```bash
  go version
  ```
- **[PostgreSQL](https://www.postgresql.org/download/)** 13+ running locally, in Docker, or hosted remotely
- **[swag CLI](https://github.com/swaggo/swag)** (only needed if you plan to regenerate the Swagger docs):
  ```bash
  go install github.com/swaggo/swag/cmd/swag@latest
  ```

## Installation

**1. Clone the repository**

```bash
git clone https://github.com/RyannWilcox/InventoryManagementAPI.git
cd inv-backend
```

**2. Install Go dependencies**

```bash
go mod download
```

**3. Set up PostgreSQL**

If you don't already have a database, create one. For example, using the `psql` CLI:

```bash
createdb inv_db
```

**4. Configure environment variables**

Create a `.env` file in the project root:

```bash
touch .env
```

Populate it with your database connection details

> The server reads these values at startup to build its PostgreSQL connection string. If `.env` is missing or a value is blank, the server will fail to start with a connection error.

## Running the Server

Start the API with:

```bash
go run main.go
```

On startup, the server will:

1. Connect to PostgreSQL using your `.env` configuration
2. Automatically migrate the `Item` schema (creating the `items` table if it doesn't exist)
3. Seed the database with 20 sample items, only if the table is currently empty
4. Start listening on `http://localhost:8080`

You should see log output confirming these steps, ending with the Gin server startup banner.

Once the server is running, Swagger documentation is available at:

```
http://localhost:8080/swagger/index.html
```

This lets you browse every endpoint, view expected request/response schemas, and send test requests directly from the browser.

## License

This project is licensed under the MIT License — see the [LICENSE](./LICENSE) file for details.
