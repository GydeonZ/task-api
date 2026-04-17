# Architecture Guide

This project follows professional Go development patterns and best practices. Here's a deep dive into the architecture.

## Design Pattern: Layered Architecture

The application is organized into distinct layers:

```
┌─────────────────────────────────────────┐
│         HTTP Handlers (HTTP Layer)      │
│    - Request parsing                    │
│    - Response formatting                │
└────────────────┬────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────┐
│    Services (Business Logic Layer)      │
│    - Application logic                  │
│    - Validation & rules                 │
│    - Transaction management             │
└────────────────┬────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────┐
│   Repositories (Data Access Layer)      │
│    - Query building                     │
│    - Result mapping                     │
│    - Error handling                     │
└────────────────┬────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────┐
│      Database (PostgreSQL)              │
└─────────────────────────────────────────┘
```

## Key Principles

### 1. Separation of Concerns
Each layer has a specific responsibility:
- **Handlers**: HTTP interface (REST endpoints)
- **Services**: Business logic and rules
- **Repositories**: Database operations

### 2. Dependency Injection
Dependencies are passed as parameters instead of created inside functions:

```go
// ✅ Good - Dependencies injected
type BoardService struct {
    repo *repository.BoardRepository
}

func NewBoardService(repo *repository.BoardRepository) *BoardService {
    return &BoardService{repo: repo}
}

// ❌ Bad - Creating dependencies inside
type BoardService struct {}
func (s *BoardService) CreateBoard() {
    repo := repository.NewBoardRepository() // Wrong!
}
```

### 3. Error Handling
Custom error types provide context:

```go
// Centralized error handling in pkg/errors/
var (
    ErrBoardNotFound = &AppError{
        Code:    "BOARD_NOT_FOUND",
        Message: "Board not found",
        Status:  404,
    }
)
```

### 4. Configuration Management
Configuration is loaded once at startup:

```go
cfg, err := config.Load()  // Loads from .env file
db.Connect(&cfg.Database)
```

## File Organization

### `cmd/server/main.go`
**Purpose**: Application entry point
- Loads configuration
- Connects to database
- Sets up Fiber app
- Initializes routes

```go
package main
func main() {
    cfg, _ := config.Load()
    database.Connect(&cfg.Database)
    app := fiber.New()
    // ... setup
    app.Listen(...)
}
```

### `internal/config/config.go`
**Purpose**: Configuration management
- Reads environment variables
- Provides defaults
- Validates configuration

```go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
}

func Load() (*Config, error) {
    // Load from .env
}
```

### `internal/models/models.go`
**Purpose**: Data structures
- Domain models (User, Board, List, Task)
- Request DTOs (CreateBoardRequest)
- Response wrappers (APIResponse)

```go
type Board struct {
    ID        int       `json:"id" db:"id"`
    UserID    int       `json:"user_id" db:"user_id"`
    Title     string    `json:"title" db:"title"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
```

### `internal/handlers/*.go`
**Purpose**: HTTP request handlers
- Parse request body
- Validate input
- Call services
- Return responses

```go
func (h *BoardHandler) CreateBoard(c *fiber.Ctx) error {
    var req models.CreateBoardRequest
    c.BodyParser(&req)
    board, _ := h.service.CreateBoard(c.Context(), userID, &req)
    return c.Status(201).JSON(board)
}
```

### `internal/service/*.go`
**Purpose**: Business logic
- Implement use cases
- Validate business rules
- Coordinate repositories
- Handle transactions

```go
func (s *BoardService) CreateBoard(ctx context.Context, userID int, req *models.CreateBoardRequest) {
    if req.Title == "" {
        return nil, errors.ErrInvalidInput
    }
    board := &models.Board{
        UserID: userID,
        Title: req.Title,
        Slug: slugify(req.Title),
    }
    return s.repo.Create(ctx, board)
}
```

### `internal/repository/*.go`
**Purpose**: Data access
- Build SQL queries
- Execute database operations
- Map results to models
- Return consistent errors

```go
func (r *BoardRepository) Create(ctx context.Context, board *models.Board) error {
    query := "INSERT INTO boards ..."
    err := database.DB.QueryRowContext(ctx, query, ...).Scan(...)
    return err
}
```

### `internal/middleware/middleware.go`
**Purpose**: Cross-cutting concerns
- Request logging
- Error recovery
- CORS handling
- Authentication (can be added)

```go
func SetupMiddleware(app *fiber.App) {
    app.Use(recover.New())                // Panic recovery
    app.Use(logger.New(...))              // Request logging
    app.Use(cors.New(cors.Config{...}))   // CORS
}
```

### `internal/database/connection.go`
**Purpose**: Database setup
- Create connection pool
- Set connection limits
- Ping database

```go
func Connect(cfg *config.DatabaseConfig) error {
    DB, _ := sql.Open("postgres", cfg.DSN)
    DB.SetMaxOpenConns(25)
    DB.SetMaxIdleConns(5)
    return DB.Ping()
}
```

### `pkg/errors/errors.go`
**Purpose**: Error definitions
- Custom error types
- Error codes for API responses
- HTTP status mapping

## Data Flow Example: Create Board

```
1. HTTP Request
   POST /api/v1/boards
   {"title": "My Board"}

2. Handler receives request
   BoardHandler.CreateBoard(c *fiber.Ctx)
   - Parse request body
   - Extract user ID from context
   
3. Call service
   boardService.CreateBoard(ctx, userID, req)
   - Validate input
   - Create Board struct
   
4. Call repository
   boardRepo.Create(ctx, board)
   - Build SQL query
   - Execute query
   - Map result to model
   
5. Return through layers
   board → service → handler → HTTP response (201 Created)
```

## Database Layer Pattern

### Raw SQL Pattern (Used in this project)
Direct SQL queries with context support:

```go
query := `
    INSERT INTO boards (user_id, title, slug, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, ...
`
err := database.DB.QueryRowContext(ctx, query, values...).Scan(&board.ID, ...)
```

**Advantages:**
- Full control over queries
- Easy to optimize
- Explicit dependencies
- Great for learning

**When to add ORM:**
- Complex queries
- Many associations
- Active record pattern needed
- Use `gorm` or `sqlc`

## Testing Patterns

### Handler Tests
```go
// Test that handlers call services correctly
func TestCreateBoard(t *testing.T) {
    // Mock service
    // Call handler
    // Verify response
}
```

### Service Tests
```go
// Test business logic
func TestBoardService_CreateBoard(t *testing.T) {
    // Mock repository
    // Call service
    // Verify logic
}
```

### Repository Tests
```go
// Test database operations
func TestBoardRepository_Create(t *testing.T) {
    // Setup test database
    // Execute operation
    // Verify database state
}
```

## Configuration Best Practices

### Environment-Specific Settings
```bash
# Development (.env)
APP_ENV=development
DB_HOST=localhost
DEBUG=true

# Production (environment variables)
APP_ENV=production
DB_HOST=prod-db.example.com
DEBUG=false
```

### Sensitive Data
**Never commit:**
- API keys in code
- Database passwords
- JWT secrets
- Credentials

**Always use:**
- `.env` files (local only)
- Environment variables (production)
- Secrets manager (AWS Secrets Manager, etc.)

## Security Considerations

### Current Implementation
- ✅ SQL injection prevention (parameterized queries)
- ✅ CORS configured
- ✅ Error messages don't expose internals

### To Add
- [ ] Authentication (JWT tokens)
- [ ] Authorization (role-based access)
- [ ] Rate limiting
- [ ] Request validation
- [ ] Password hashing
- [ ] HTTPS/TLS
- [ ] CSRF protection

## Performance Optimization

### Current Setup
- ✅ Connection pooling (MaxOpenConns: 25, MaxIdleConns: 5)
- ✅ Indexed database columns
- ✅ Raw queries (efficient)

### To Add
- [ ] Caching layer (Redis)
- [ ] Query result caching
- [ ] Pagination for list endpoints
- [ ] Batch operations
- [ ] Database query optimization
- [ ] Async background jobs

## Dependency Management

### go.mod Structure
```go
module github.com/GydeonZ/task-api

require (
    github.com/gofiber/fiber/v2 v2.52.1
    github.com/lib/pq v1.10.9  // PostgreSQL driver
    github.com/joho/godotenv v1.5.1
)
```

### Adding Dependencies
```bash
go get github.com/package/name
go mod tidy  # Remove unused dependencies
```

## Deployment Checklist

- [ ] Set production environment variables
- [ ] Use connection pooling
- [ ] Enable HTTPS
- [ ] Set up logging
- [ ] Configure monitoring
- [ ] Run database migrations
- [ ] Test error scenarios
- [ ] Set up backups
- [ ] Use secrets manager
- [ ] Configure CI/CD pipeline
