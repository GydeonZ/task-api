# Code Standards & Best Practices

Follow these standards to keep the project maintainable and professional.

## Naming Conventions

### Files
```go
// ✅ Correct: snake_case
board_handler.go
board_service.go
board_repository.go

// ❌ Wrong
BoardHandler.go
boardHandler.go
board-handler.go
```

### Functions
```go
// ✅ Correct: PascalCase for exported, camelCase for private
func (h *BoardHandler) CreateBoard(c *fiber.Ctx) error {}
func slugify(s string) string {}

// ❌ Wrong
func (h *BoardHandler) createBoard(c *fiber.Ctx) error {}
func (h *BoardHandler) create_board(c *fiber.Ctx) error {}
```

### Types & Structs
```go
// ✅ Correct: PascalCase, descriptive names
type BoardHandler struct {}
type CreateBoardRequest struct {}

// ❌ Wrong
type boardHandler struct {}
type BoardReq struct {}  // Too short
```

### Constants
```go
// ✅ Correct: SCREAMING_SNAKE_CASE
const (
    MAX_TITLE_LENGTH = 255
    DEFAULT_PAGE_SIZE = 20
)

// ❌ Wrong
const (
    max_title_length = 255
    maxTitleLength = 255
)
```

### Variables
```go
// ✅ Correct: camelCase, descriptive
userID := 1
boardTitle := "My Board"

// ❌ Wrong
u := 1           // Too short
board_title := "" // Too verbose
```

## Code Organization

### Imports
```go
// ✅ Correct: Organized by category
package handlers

import (
    "strconv"      // Standard library
    
    "github.com/gofiber/fiber/v2"  // Third party
    
    "github.com/GydeonZ/task-api/internal/models"  // Project
    "github.com/GydeonZ/task-api/internal/service"
)

// ❌ Wrong: Random order
import (
    "github.com/GydeonZ/task-api/internal/service"
    "strconv"
    "github.com/gofiber/fiber/v2"
)
```

### Function Organization
```go
// ✅ Correct: Group related functions
type BoardHandler struct { /* ... */ }

func NewBoardHandler(service *BoardService) *BoardHandler { /* ... */ }

func (h *BoardHandler) CreateBoard(c *fiber.Ctx) error { /* ... */ }
func (h *BoardHandler) GetBoard(c *fiber.Ctx) error { /* ... */ }
func (h *BoardHandler) UpdateBoard(c *fiber.Ctx) error { /* ... */ }

// ❌ Wrong: Random organization
func (h *BoardHandler) GetBoard(c *fiber.Ctx) error { /* ... */ }
func NewBoardHandler(service *BoardService) *BoardHandler { /* ... */ }
func (h *BoardHandler) CreateBoard(c *fiber.Ctx) error { /* ... */ }
```

## Error Handling

### Always Handle Errors
```go
// ✅ Correct
rows, err := database.DB.QueryContext(ctx, query)
if err != nil {
    return errors.NewWithErr("DB_ERROR", "Failed to query", 500, err)
}
defer rows.Close()

// ❌ Wrong
rows, _ := database.DB.QueryContext(ctx, query)  // Ignoring error!
defer rows.Close()
```

### Return Meaningful Errors
```go
// ✅ Correct
if req.Title == "" {
    return nil, errors.New("INVALID_INPUT", "Title cannot be empty", 400)
}

// ❌ Wrong
if req.Title == "" {
    return nil, errors.New("ERROR")  // Unclear
}
```

### Wrap Errors with Context
```go
// ✅ Correct
db := database.DB.QueryRowContext(ctx, query, args...)
if err != nil {
    return nil, errors.NewWithErr("DB_ERROR", 
        "Failed to create board", 500, err)
}

// ❌ Wrong
if err := db.Scan(...); err != nil {
    return nil, err  // Lost context
}
```

## Comments

### Comment Guidelines

Write comments for **why**, not **what**:

```go
// ✅ Correct: Explain the reasoning
// We use a connection pool to reuse connections
// and improve database performance
DB.SetMaxOpenConns(25)

// ❌ Wrong: Describing obvious code
// Set max open connections to 25
DB.SetMaxOpenConns(25)
```

### Function Comments

```go
// ✅ Correct: Export comments
// CreateBoard creates a new board for a user.
// It validates input and returns the created board or an error.
func (s *BoardService) CreateBoard(ctx context.Context, userID int, req *models.CreateBoardRequest) (*models.Board, error) {

// ❌ Wrong: Missing or unclear
func (s *BoardService) CreateBoard(ctx context.Context, userID int, req *models.CreateBoardRequest) (*models.Board, error) {
```

### Package Comments

```go
// ✅ Correct: At package level
// Package handlers provides HTTP request handlers
// for the task management API.
package handlers

// ❌ Wrong: Missing or vague
package handlers
```

## Testing

### Test Naming
```go
// ✅ Correct: Clear what's being tested
func TestBoardService_CreateBoard_Success(t *testing.T) {}
func TestBoardService_CreateBoard_InvalidInput(t *testing.T) {}
func TestBoardRepository_Create_DatabaseError(t *testing.T) {}

// ❌ Wrong
func TestCreate(t *testing.T) {}
func Test1(t *testing.T) {}
```

### Table-Driven Tests
```go
// ✅ Correct: Test multiple scenarios
func TestValidateTitle(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"Valid", "My Board", false},
        {"Empty", "", true},
        {"Long", strings.Repeat("a", 256), true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateTitle(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("unexpected error: %v", err)
            }
        })
    }
}
```

### Use Asserts
```bash
# Install assertion library
go get github.com/stretchr/testify/assert
```

```go
// ✅ Correct: Clear assertions
assert.NoError(t, err)
assert.Equal(t, expected, actual)
assert.NotNil(t, result)

// ❌ Wrong: Manual checks
if err != nil {
    t.Fatalf("unexpected error: %v", err)
}
```

## Database Patterns

### Use Context
```go
// ✅ Correct: Always use context
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

row := database.DB.QueryRowContext(ctx, query, args...)

// ❌ Wrong: No context
row := database.DB.QueryRow(query, args...)
```

### Always Close Resources
```go
// ✅ Correct: Defer close immediately
rows, err := database.DB.QueryContext(ctx, query)
if err != nil {
    return err
}
defer rows.Close()

// ❌ Wrong: Risk of leaked connections
rows, _ := database.DB.QueryContext(ctx, query)
// ... lots of code ...
rows.Close()
```

### Use Transactions for Related Operations
```go
// ✅ Correct: Transactional integrity
tx, err := database.DB.BeginTx(ctx, nil)
if err != nil {
    return err
}
defer tx.Rollback()

// Multiple operations
_, err = tx.ExecContext(ctx, "INSERT INTO boards ...")
if err != nil {
    return err
}

_, err = tx.ExecContext(ctx, "INSERT INTO lists ...")
if err != nil {
    return err
}

return tx.Commit().Error
```

## Middleware & Handlers

### Handler Structure
```go
// ✅ Correct: Clear, simple handlers
func (h *BoardHandler) CreateBoard(c *fiber.Ctx) error {
    // 1. Parse request
    var req models.CreateBoardRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(models.APIResponse{
            Success: false,
            Error:   "Invalid request",
        })
    }
    
    // 2. Call service
    board, err := h.service.CreateBoard(c.Context(), userID, &req)
    if err != nil {
        return c.Status(500).JSON(models.APIResponse{
            Success: false,
            Error:   err.Error(),
        })
    }
    
    // 3. Return response
    return c.Status(201).JSON(models.APIResponse{
        Success: true,
        Data:    board,
    })
}

// ❌ Wrong: Mixing concerns
func (h *BoardHandler) CreateBoard(c *fiber.Ctx) error {
    var req models.CreateBoardRequest
    c.BodyParser(&req)
    
    // Business logic in handler (WRONG!)
    if req.Title == "" {
        return errors.New("invalid")
    }
    
    board := &models.Board{Title: req.Title}
    database.DB.Exec("INSERT INTO boards ...")  // DB query in handler!
    
    return c.JSON(board)
}
```

## Configuration

### Use Constants for Magic Values
```go
// ✅ Correct
const (
    MAX_TITLE_LENGTH = 255
    MAX_OPEN_CONNS = 25
    MAX_IDLE_CONNS = 5
    REQUEST_TIMEOUT = 30 * time.Second
)

DB.SetMaxOpenConns(MAX_OPEN_CONNS)

// ❌ Wrong
DB.SetMaxOpenConns(25)  // What is 25?
if len(title) > 255 {   // Why 255?
}
```

### Environment Variables in Config
```go
// ✅ Correct: Config from environment
type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
}

cfg.Database.Host = getEnv("DB_HOST", "localhost")

// ❌ Wrong: Hardcoded values
const DB_HOST = "localhost"  // Can't change without recompile
```

## Performance

### Use Prepared Statements for Repeated Queries
```go
// ✅ Correct: Prepare once, use many times
stmt, err := database.DB.PrepareContext(ctx, query)
defer stmt.Close()

for _, board := range boards {
    stmt.ExecContext(ctx, board.ID, board.Title)
}

// ❌ Wrong: Prepare every time
for _, board := range boards {
    database.DB.Exec(query, board.ID, board.Title)
}
```

### Batch Operations
```go
// ✅ Correct: Combine multiple operations
query := `
    INSERT INTO tasks (list_id, title) VALUES
    ($1, $2),
    ($3, $4),
    ($5, $6)
`
database.DB.ExecContext(ctx, query, args...)

// ❌ Wrong: One-by-one inserts
for _, task := range tasks {
    database.DB.Exec("INSERT INTO tasks ...")
}
```

## Logging Strategy

### Log Levels

```go
// INFO: Important business events
log.Info("User registered", "user_id", userID)
log.Info("Board created", "board_id", boardID)

// ERROR: Errors that need attention
log.Error("Database connection failed", "error", err)
log.Error("User not found", "user_id", userID)

// DEBUG: Detailed information for debugging
log.Debug("SQL query executed", "query", query, "args", args)
log.Debug("Request headers", "headers", c.GetReqHeaders())

// WARN: Warning-level events
log.Warn("Slow query", "duration", 2*time.Second)
log.Warn("High memory usage", "percent", 85)
```

## Security Checklist

- [ ] No hardcoded passwords or secrets
- [ ] Use parameterized SQL queries (prevent SQL injection)
- [ ] Validate and sanitize all user input
- [ ] Use HTTPS in production
- [ ] Hash passwords (never store plaintext)
- [ ] Implement rate limiting
- [ ] Add CORS headers appropriately
- [ ] Use environment variables for secrets
- [ ] Implement proper authentication
- [ ] Log security events

## Code Review Checklist

Before submitting code:

- [ ] Follows naming conventions
- [ ] Error handling is proper
- [ ] No hardcoded values
- [ ] Tests added for new features
- [ ] Comments explain "why" not "what"
- [ ] No security vulnerabilities
- [ ] Performance acceptable
- [ ] Follows project structure
- [ ] Database connections properly closed
- [ ] Context timeouts implemented

## Common Mistakes to Avoid

❌ **Don't:**
- Ignore errors with `_`
- Use global variables for state
- Create deeply nested functions
- Mix concerns in one layer
- Hardcode configuration
- Forget to close resources
- Use string concatenation in SQL (SQL injection!)
- Commit sensitive data
- Leave console output in production code
-Use panic for error handling (except at startup)
