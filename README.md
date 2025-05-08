# warehouse-api

This is a simple Warehouse CRUD (Create, Read, Update, Delete) application built with **Go** using the following technologies:
- **Gin** for HTTP routing
- **SQLC** with PostgreSQL

## Features
- Basic CRUD operations for products
- Graceful error handling with appropriate HTTP status codes
- Auto-generated API documentation (Swagger)

## Prerequisites
- Go 1.20+
- Sqlc & migrate programs installed
- PostgreSQL:15.8 container in Docker

## Installation & Setup
1. **Clone the repository**:
```bash
git clone https://github.com/Gen1usBruh/warehouse-api.git
cd warehouse-api
```
2. **Create docker container for postgres**:
```bash
sudo docker run --name pg_warehouse -p 5433:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d postgres_image_id
```

3. **Install dependencies**:
```bash
go mod tidy
```

4. **Run the application locally**:
```bash
go run main.go
```

## API Endpoints
- `GET /products` - Get list of all products.
- `DELETE /products/:id` - Delete a product by id.
- `PUT /products/:id` - Update a product by id.
- `GET /products/:id` - Get a product by id.
- `POST /products` - Add a new product.
