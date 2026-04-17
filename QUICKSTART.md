# Quick Start Guide

## Setup & Installation

### 1. Install Go Dependencies
```bash
go mod download
```

### 2. Setup Environment Variables
```bash
cp .env.example .env
```

### 3. Start PostgreSQL Database
```bash
# Option 1: Using Docker
docker run -d --name task_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=task_management \
  -p 5432:5432 \
  postgres:16-alpine

# Option 2: Using Docker Compose
docker-compose up -d
```

### 4. Create Database Tables (Migrations)
```bash
# Download psql if not installed, then:
psql -h localhost -U postgres -d task_management -f migrations/001_initial_schema.sql
```

### 5. Run the Application
```bash
# Development mode (with hot reload, if using air)
go run ./cmd/server/main.go

# Or using Make
make dev

# Or build and run
make build
make run
```

The API will be available at `http://localhost:3000`

## API Usage Examples

### Health Check
```bash
curl http://localhost:3000/health
```

### Create a Board
```bash
curl -X POST http://localhost:3000/api/v1/boards \
  -H "Content-Type: application/json" \
  -d '{"title": "My First Board"}'
```

### Get All Boards
```bash
curl http://localhost:3000/api/v1/boards
```

### Get a Specific Board
```bash
curl http://localhost:3000/api/v1/boards/1
```

### Create a List in a Board
```bash
curl -X POST http://localhost:3000/api/v1/boards/1/lists \
  -H "Content-Type: application/json" \
  -d '{"title": "To Do"}'
```

### Get All Lists in a Board
```bash
curl http://localhost:3000/api/v1/boards/1/lists
```

### Create a Task in a List
```bash
curl -X POST http://localhost:3000/api/v1/lists/1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Fix login bug",
    "description": "Users unable to login with special characters",
    "due_date": "2026-05-10T00:00:00Z"
  }'
```

### Get All Tasks in a List
```bash
curl http://localhost:3000/api/v1/lists/1/tasks
```

### Update a Task
```bash
curl -X PUT http://localhost:3000/api/v1/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Fix login bug",
    "description": "Updated description",
    "is_completed": false
  }'
```

### Mark Task as Completed
```bash
curl -X PATCH http://localhost:3000/api/v1/tasks/1/complete
```

### Delete a Task
```bash
curl -X DELETE http://localhost:3000/api/v1/tasks/1
```

## Project Structure Explanation

```
task-api/
├── cmd/server/main.go              # Application entry point
├── internal/
│   ├── config/config.go            # Load & manage configuration
│   ├── database/connection.go       # PostgreSQL connection & setup
│   ├── models/models.go            # Data models & DTOs
│   ├── handlers/                   # HTTP request handlers
│   │   ├── board_handler.go
│   │   ├── list_handler.go
│   │   └── task_handler.go
│   ├── routes/routes.go            # API route definitions
│   ├── middleware/middleware.go    # CORS, logger, recover
│   ├── repository/                 # Data access layer
│   │   ├── board_repository.go
│   │   ├── list_repository.go
│   │   └── task_repository.go
│   └── service/                    # Business logic layer
│       ├── board_service.go
│       ├── list_service.go
│       └── task_service.go
├── pkg/errors/errors.go            # Error handling
├── migrations/                     # Database schemas
└── docker-compose.yml              # Local development setup
```

## Architecture Layers

### 1. **Handlers** (HTTP Layer)
- Parses incoming requests
- Calls services with validated data
- Returns JSON responses

### 2. **Services** (Business Logic)
- Contains application logic
- Validates business rules
- Calls repositories

### 3. **Repositories** (Data Access)
- Abstracts database operations
- Executes SQL queries
- Handles data errors

## Next Steps

1. **Add Authentication**
   - Implement JWT token generation
   - Add user login/registration endpoints
   - Protect endpoints with middleware

2. **Add Validation**
   - Use `go-playground/validator` for request validation
   - Add custom validation rules

3. **Add Tests**
   - Unit tests for services
   - Integration tests for handlers
   - Test database migrations

4. **Add Logging**
   - Use `logrus` or `zap` for structured logging
   - Log all database operations

5. **Add API Documentation**
   - Generate Swagger/OpenAPI docs
   - Add interactive API explorer

6. **Performance**
   - Add caching layer (Redis)
   - Add database indexing
   - Implement pagination

7. **Deployment**
   - Configure CI/CD pipeline
   - Set up production environment variables
   - Monitor application health
