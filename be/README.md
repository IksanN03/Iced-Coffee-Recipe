# Iced Coffee Recipe Management API

A backend service for managing iced coffee recipes and inventory built with Go, Gin, and PostgreSQL.

## Prerequisites
- Go 1.22 or higher
- PostgreSQL 13
- Docker & Docker Compose

## Step
- Clone the repository: git clone <repository-url>
- Rename the .env.example file to .env and update the environment variables.
- Start the frontend application first and note its URL. Then update the FRONTEND_URL in your .env file with that URL.
- Run 'docker-compose up --build' OR 'go run main.go' to start the application and database.
- For running with cmd 'go run main.go': set DB_HOST=localhost
- For running with cmd 'docker-compose up --build': set DB_HOST=db
- Access the API at http://localhost:${PORT}.

## API Endpoints
- Authentication
POST /auth/submit-email - Request magic link
GET /auth/magic-link - Verify magic link

- Inventory Management
GET /inventory - List inventory items
POST /inventory - Add new item
PUT /inventory/:id - Update item
DELETE /inventory/:id - Delete item

- Recipe Management
GET /recipe - List all recipes
POST /recipe - Create new recipe
PUT /recipe/:id - Update recipe

## Database Schema
The service uses PostgreSQL with the following main tables:
- users
- inventory
- recipes

## Test Database Setup (test.sql)
Download the test.sql file to set up your test database. This file contains:
- Table creation scripts
- Sample data for inventory items
- Sample recipe data
- Test user accounts

To use the test database:
```bash
psql -U your_username -d your_database -f test.sql
```

## API Testing (test.postman_collection.json)
A complete Postman collection is provided for testing all API endpoints.

To use the Postman collection:

- Download test.postman_collection.json
- Open Postman
- Click "Import" -> "Upload Files"
- Select test.postman_collection.json
- The collection includes:
- Environment variables setup
- Authentication requests
- Inventory management requests
- Recipe management requests
- Test data examples
- Pre-request scripts
- Response validation tests

## Testing
Run the test suite:

go test ./...