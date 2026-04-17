# Troubleshooting Guide

## Common Issues and Solutions

### 1. Database Connection Failed

**Error**: `failed to connect to database`

**Causes:**
- PostgreSQL not running
- Wrong credentials in `.env`
- Database doesn't exist
- Wrong host/port

**Solutions:**
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Verify connection details in .env
cat .env

# Try connecting directly
psql -h localhost -U postgres -d task_management

# If using Docker Compose
docker-compose logs postgres

# Restart database
docker-compose restart postgres
```

### 2. Port Already in Use

**Error**: `listen tcp :3000: bind: address already in use`

**Solution:**
```bash
# Find process on port 3000
lsof -i :3000  # macOS/Linux

# Kill the process
kill -9 <PID>

# Or change port in .env
SERVER_PORT=3001
```

### 3. Migration Failed

**Error**: `migration failed` or `table already exists`

**Solutions:**
```bash
# Connect to database
psql -h localhost -U postgres -d task_management

# Drop all tables and start fresh
DROP TABLE tasks CASCADE;
DROP TABLE lists CASCADE;
DROP TABLE boards CASCADE;
DROP TABLE users CASCADE;

# Re-run migration
psql -h localhost -U postgres -d task_management -f migrations/001_initial_schema.sql

# Or delete database and recreate
docker-compose restart postgres
psql -h localhost -U postgres -d task_management -f migrations/001_initial_schema.sql
```

### 4. Module Import Errors

**Error**: `no required module providing package`

**Solutions:**
```bash
# Download dependencies
go mod download

# Clean and rebuild
go clean -modcache
go mod tidy
go build ./cmd/server

# Check go.mod for inconsistencies
cat go.mod
```

### 5. Build Fails

**Error**: `undefined reference` or compile errors

**Solutions:**
```bash
# Ensure you're in project root
pwd
cat go.mod

# Clean build
go clean
go build ./cmd/server

# Check for syntax errors
go fmt ./...
go vet ./...

# Build with verbose output
go build -v ./cmd/server
```

### 6. Docker Build Issues

**Error**: `failed to build Docker image`

**Solutions:**
```bash
# Check Dockerfile syntax
docker build -t task-api .

# Check for missing files
ls -la migrations/
ls -la cmd/server/

# Build without cache
docker build --no-cache -t task-api .

# View build logs
docker build -t task-api . 2>&1 | tee build.log
```

### 7. CORS Errors in Frontend

**Error**: `Access to XMLHttpRequest blocked by CORS policy`

**Solution** (in `internal/middleware/middleware.go`):
```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "http://localhost:3000,http://localhost:3001",  // Add your frontend URL
    AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
    AllowHeaders: "Origin,Content-Type,Accept,Authorization",
}))
```

### 8. JSON Parsing Errors

**Error**: `failed to parse request body` or `invalid JSON`

**Solutions:**
```bash
# Check Content-Type header
curl -X POST http://localhost:3000/api/v1/boards \
  -H "Content-Type: application/json" \
  -d '{"title": "Test"}'

# Check JSON validity
echo '{"title": "Test"}' | jq .

# Verify request format
# ✅ Correct
{"title": "My Board"}

# ❌ Wrong (unquoted strings, trailing commas, etc.)
{title: "My Board",}
```

### 9. 404 Not Found on Endpoints

**Error**: `{"success": false, "error": "Route not found"}`

**Causes:**
- Wrong URL path
- Wrong HTTP method
- Route not registered

**Solutions:**
```bash
# Verify endpoint exists in routes.go
cat internal/routes/routes.go

# Test with curl
curl -X GET http://localhost:3000/api/v1/boards
curl -X POST http://localhost:3000/api/v1/boards

# Check request URL exactly
# ✅ Correct
POST /api/v1/boards

# ❌ Wrong
POST /api/boards  # Missing /v1
POST /API/V1/boards  # Case sensitive
POST /api/v1/boards/ # Extra slash
```

### 10. Database Locked Error

**Error**: `database is locked` or `connection refused`

**Solutions:**
```bash
# Check active connections
docker-compose logs postgres

# Kill existing connections
psql -h localhost -U postgres <<EOF
SELECT pg_terminate_backend(pid) FROM pg_stat_activity 
WHERE datname = 'task_management';
EOF

# Restart PostgreSQL
docker-compose restart postgres
```

## Debug Tips

### Enable Verbose Logging

```go
// In internal/middleware/middleware.go
app.Use(logger.New(logger.Config{
    Format: "[${ip}] ${status} - ${method} ${path} | ${latency}\n",
}))
```

### Print SQL Queries

```go
// In repositories
log.Printf("SQL: %s", query)
log.Printf("Args: %v", args)
```

### Test Database Directly

```bash
# Connect to database
psql -h localhost -U postgres -d task_management

# List tables
\dt

# View table structure
\d boards

# Sample query
SELECT * FROM boards LIMIT 10;

# Check row count
SELECT COUNT(*) FROM tasks;
```

### Check API Responses

```bash
# With headers
curl -i http://localhost:3000/api/v1/boards

# Pretty JSON
curl http://localhost:3000/api/v1/boards | jq .

# With error info
curl -v http://localhost:3000/api/v1/boards 2>&1 | grep '>'
```

## Performance Issues

### Slow Queries

**Check indexes:**
```sql
-- From migrations/001_initial_schema.sql
-- Verify indexes exist
SELECT * FROM pg_indexes WHERE tablename = 'tasks';
```

**Add missing indexes:**
```sql
CREATE INDEX idx_tasks_created_at ON tasks(created_at);
CREATE INDEX idx_boards_created_at ON boards(created_at);
```

### High Memory Usage

**Solutions:**
```go
// In internal/database/connection.go
DB.SetMaxOpenConns(10)    // Reduce from 25
DB.SetMaxIdleConns(2)     // Reduce from 5

// Or in docker-compose.yml
environment:
  POSTGRES_SHARED_BUFFERS: 256MB
```

### Connection Pool Exhaustion

**Symptoms:**
- "too many connections" errors
- Slow responses
- Increased latency

**Solutions:**
```go
// Ensure connections are closed
defer rows.Close()
defer result.Close()

// Use context timeouts
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

## Getting Help

### Check Logs
```bash
# Application logs
docker-compose logs task_api

# Database logs
docker-compose logs postgres

# Follow logs in real-time
docker-compose logs -f api
```

### Useful Commands

```bash
# Health check
curl http://localhost:3000/health

# Test specific endpoint
curl http://localhost:3000/api/v1/boards

# Database stats
docker exec task_api_postgres psql -U postgres -d task_management -c "SELECT * FROM pg_stat_statements LIMIT 10;"

# Docker stats
docker stats task_api_server task_api_postgres
```

### Clean Slate

```bash
# Stop everything
docker-compose down

# Remove volumes (WARNING: deletes data!)
docker-compose down -v

# Start fresh
docker-compose up -d

# Rebuild containers
docker-compose build --no-cache
docker-compose up -d
```
