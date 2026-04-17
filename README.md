# Task Management API

A professional task management system (Trello-like) built with Go, Fiber, and PostgreSQL.

## Features

- ✅ Create and manage boards
- ✅ Create lists within boards
- ✅ Create tasks within lists
- ✅ Mark tasks as completed
- ✅ RESTful API
- ✅ PostgreSQL database
- ✅ Docker support
- ✅ Error handling

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/                     # Configuration management
│   ├── database/                   # Database connection
│   ├── models/                     # Data models and DTOs
│   ├── handlers/                   # HTTP handlers
│   ├── routes/                     # API routes
│   ├── middleware/                 # Middleware
│   ├── repository/                 # Data access layer
│   └── service/                    # Business logic
├── pkg/
│   └── errors/                     # Error handling
├── migrations/                     # Database migrations
├── go.mod                          # Go module file
├── go.sum                          # Go module checksum
├── Dockerfile                      # Docker image
├── docker-compose.yml              # Docker Compose configuration
├── Makefile                        # Build commands
└── .env.example                    # Environment variables template
```

## Prerequisites

- Go 1.22+
- PostgreSQL 16+
- Docker & Docker Compose (optional)
- Make (optional)

## Installation

1. Clone the repository:
```bash
git clone <repository>
cd task-api
```

2. Copy environment variables:
```bash
cp .env.example .env
```

3. Install dependencies:
```bash
go mod download
```

## Running the Application

### Local Development

1. Start PostgreSQL (using Docker):
```bash
docker run -d --name task_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=task_management \
  -p 5432:5432 \
  postgres:16-alpine
```

2. Run migrations:
```bash
psql -h localhost -U postgres -d task_management -f migrations/001_initial_schema.sql
```

3. Run the application:
```bash
go run ./cmd/server/main.go
```

Or using Make:
```bash
make dev
```

### Using Docker Compose

```bash
docker-compose up -d
```

This will start both PostgreSQL and the API server.

## API Endpoints

### Boards
- `GET /api/v1/boards` - Get all user boards
- `GET /api/v1/boards/:id` - Get a specific board
- `POST /api/v1/boards` - Create a new board
- `PUT /api/v1/boards/:id` - Update a board
- `DELETE /api/v1/boards/:id` - Delete a board

### Lists
- `GET /api/v1/boards/:boardId/lists` - Get all lists in a board
- `GET /api/v1/lists/:id` - Get a specific list
- `POST /api/v1/boards/:boardId/lists` - Create a new list
- `PUT /api/v1/lists/:id` - Update a list
- `DELETE /api/v1/lists/:id` - Delete a list

### Tasks
- `GET /api/v1/lists/:listId/tasks` - Get all tasks in a list
- `GET /api/v1/tasks/:id` - Get a specific task
- `POST /api/v1/lists/:listId/tasks` - Create a new task
- `PUT /api/v1/tasks/:id` - Update a task
- `DELETE /api/v1/tasks/:id` - Delete a task
- `PATCH /api/v1/tasks/:id/complete` - Mark task as completed

## Example Requests

### Create a Board
```bash
curl -X POST http://localhost:3000/api/v1/boards \
  -H "Content-Type: application/json" \
  -d '{"title": "My Project"}'
```

### Create a List
```bash
curl -X POST http://localhost:3000/api/v1/boards/1/lists \
  -H "Content-Type: application/json" \
  -d '{"title": "To Do"}'
```

### Create a Task
```bash
curl -X POST http://localhost:3000/api/v1/lists/1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Implement API", "description": "Build REST endpoints"}'
```

## Make Commands

```bash
make build       # Build the application
make run         # Build and run
make dev         # Run in development mode
make test        # Run tests
make clean       # Clean build artifacts
make db-up       # Start PostgreSQL
make db-down     # Stop PostgreSQL
make docker-up   # Start services with Docker
make docker-down # Stop services
make migrate     # Run migrations
make fmt         # Format code
make lint        # Run linter
make deps        # Download dependencies
```

## Database Schema

### Users
- `id` - Primary key
- `username` - Unique username
- `email` - Unique email
- `password_hash` - Hashed password
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### Boards
- `id` - Primary key
- `user_id` - Foreign key to users
- `title` - Board title
- `slug` - URL-friendly title
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### Lists
- `id` - Primary key
- `board_id` - Foreign key to boards
- `title` - List title
- `position` - Display order
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### Tasks
- `id` - Primary key
- `list_id` - Foreign key to lists
- `title` - Task title
- `description` - Task description
- `position` - Display order
- `is_completed` - Completion status
- `due_date` - Optional due date
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

## Architecture Principles

This project follows senior Go developer best practices:

1. **Layer Architecture**: Clear separation of concerns (handlers, services, repositories)
2. **Dependency Injection**: Services receive dependencies as parameters
3. **Error Handling**: Centralized error types and handling
4. **Configuration Management**: Environment-based configuration
5. **Database Connection Pool**: Optimized connection management
6. **Middleware Pattern**: Reusable middleware for cross-cutting concerns
7. **Repository Pattern**: Data access abstraction
8. **Service Layer**: Business logic separation
9. **Standard Project Layout**: Following Go project structure conventions

## Future Enhancements

- [ ] Authentication & Authorization (JWT)
- [ ] User validation
- [ ] Input validation with decorators
- [ ] Test suite
- [ ] API documentation (Swagger)
- [ ] Logging system
- [ ] Caching layer
- [ ] Rate limiting
- [ ] WebSockets for real-time updates
- [ ] File upload support
- [ ] Search functionality
- [ ] Activity tracking

## License

This project is licensed under the MIT License - see LICENSE file for details.