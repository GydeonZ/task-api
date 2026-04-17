# API Reference & Visual Guide

## Request/Response Flow

```
CLIENT
  │
  └─► HTTP REQUEST ──┐
                      │
                    Router
                      │
              ┌───────┼───────┐
              │       │       │
          Middleware  Handler
              │       │
              └───┬───┘
                  │
              Validation
                  │
              ┌───────────────────┬───────┐
              │                   │       │
            Service         Slug Generation
              │                   │
          Business Logic ◄────────┘
              │
              ├─ Validate input
              ├─ Transform data
              └─ Call repository
                  │
            Repository
              │
        ┌─────┴─────┐
        │            │
    Build SQL    Database
        │            │
        └────┬───────┘
             │
    Error Handling
      │
      └─► APIResponse
             │
             └─► HTTP Response 
                  (JSON)
                  │
                  └─► CLIENT
```

## Complete Request Example

### POST /api/v1/login - LOGIN USER

**Request:**
```bash
curl -X POST http://localhost:3000/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

**Request Flow:**

```
1. HTTP Request reaches Fiber Router
   → POST /api/v1/login

2. Route Handler matches
   → UserHandler.Login()

3. Parse Request Body
   {
     "email": "user@example.com",
     "password": "password123"
   }

4. Service Layer (Business Logic)
   → UserService.LoginUser()
   → Find user by email in database
   → Compare password with bcrypt hash
   → Generate JWT token with claims

5. Response Built
   {
     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE3MjYwMDAwMDB9.signature"
   }

6. HTTP Response Sent
   → Status: 200 OK
   → Body: JSON with token
```

### POST /api/v1/boards - CREATE BOARD

**Request:**
```bash
curl -X POST http://localhost:3000/api/v1/boards \
  -H "Content-Type: application/json" \
  -d '{"title": "My First Project"}'
```

**Request Flow:**

```
1. HTTP Request reaches Fiber Router
   → POST /api/v1/boards

2. Route Handler matches
   → BoardHandler.CreateBoard()

3. Parse Request Body
   {
     "title": "My First Project"
   }

4. ServiceLayer (Business Logic)
   → BoardService.CreateBoard()
   → Validate: Title not empty
   → Generate: Slug = "my-first-project"
   → Create: Board struct

5. Repository Layer (Data Access)
   → BoardRepository.Create()
   → Build SQL: INSERT INTO boards ...
   → Execute: database.DB.QueryRowContext()
   → Scan result: board.ID = 1

6. Response Built
   {
     "success": true,
     "message": "Board created successfully",
     "data": {
       "id": 1,
       "user_id": 1,
       "title": "My First Project",
       "slug": "my-first-project",
       "created_at": "2026-04-09T10:30:00Z",
       "updated_at": "2026-04-09T10:30:00Z"
     }
   }

7. HTTP Response Sent
   → Status: 201 Created
   → Body: JSON response
```

## All Endpoints Reference

### Health Check
```
GET /health

Response 200:
{
  "status": "ok"
}
```

### USER/AUTHENTICATION

#### Login
```
POST /api/v1/login

Request:
{
  "email": "user@example.com",
  "password": "password123"
}

Response 200:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}

Response 401:
{
  "message": "invalid credentials"
}
```

### BOARDS

#### Create Board
```
POST /api/v1/boards

Request:
{
  "title": "My Board"
}

Response 201:
{
  "success": true,
  "message": "Board created successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "title": "My Board",
    "slug": "my-board",
    "created_at": "2026-04-09T...",
    "updated_at": "2026-04-09T..."
  }
}
```

#### Get User's Boards
```
GET /api/v1/boards

Response 200:
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "title": "My Board",
      "slug": "my-board",
      ...
    }
  ]
}
```

#### Get Specific Board
```
GET /api/v1/boards/1

Response 200:
{
  "success": true,
  "data": { ... }
}

Response 404:
{
  "success": false,
  "error": "Board not found"
}
```

#### Update Board
```
PUT /api/v1/boards/1

Request:
{
  "title": "Updated Title"
}

Response 200:
{
  "success": true,
  "message": "Board updated successfully",
  "data": { ... }
}
```

#### Delete Board
```
DELETE /api/v1/boards/1

Response 200:
{
  "success": true,
  "message": "Board deleted successfully"
}
```

### LISTS

#### Create List
```
POST /api/v1/boards/1/lists

Request:
{
  "title": "To Do"
}

Response 201:
{
  "success": true,
  "message": "List created successfully",
  "data": {
    "id": 1,
    "board_id": 1,
    "title": "To Do",
    "position": 1,
    ...
  }
}
```

#### Get Board's Lists
```
GET /api/v1/boards/1/lists

Response 200:
{
  "success": true,
  "data": [ ... ]
}
```

#### Get Specific List
```
GET /api/v1/lists/1

Response: { ... }
```

#### Update List
```
PUT /api/v1/lists/1

Request:
{
  "title": "Updated List",
  "position": 2
}

Response 200:
{
  "success": true,
  "message": "List updated successfully",
  "data": { ... }
}
```

#### Delete List
```
DELETE /api/v1/lists/1

Response 200:
{
  "success": true,
  "message": "List deleted successfully"
}
```

### TASKS

#### Create Task
```
POST /api/v1/lists/1/tasks

Request:
{
  "title": "Fix login bug",
  "description": "Users can't login with special characters",
  "due_date": "2026-05-10T00:00:00Z"
}

Response 201:
{
  "success": true,
  "message": "Task created successfully",
  "data": {
    "id": 1,
    "list_id": 1,
    "title": "Fix login bug",
    "description": "...",
    "position": 1,
    "is_completed": false,
    "due_date": "2026-05-10T00:00:00Z",
    ...
  }
}
```

#### Get List's Tasks
```
GET /api/v1/lists/1/tasks

Response 200:
{
  "success": true,
  "data": [ ... ]
}
```

#### Get Specific Task
```
GET /api/v1/tasks/1

Response: { ... }
```

#### Update Task
```
PUT /api/v1/tasks/1

Request:
{
  "title": "Fix login bug",
  "description": "Updated description",
  "is_completed": false,
  "due_date": "2026-05-15T00:00:00Z",
  "position": 2
}

Response 200:
{
  "success": true,
  "message": "Task updated successfully",
  "data": { ... }
}
```

#### Mark Task Completed
```
PATCH /api/v1/tasks/1/complete

Response 200:
{
  "success": true,
  "message": "Task marked as completed",
  "data": {
    ...
    "is_completed": true,
    ...
  }
}
```

#### Delete Task
```
DELETE /api/v1/tasks/1

Response 200:
{
  "success": true,
  "message": "Task deleted successfully"
}
```

## Response Formats

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { /* entity */ }
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error message"
}
```

### List Response
```json
{
  "success": true,
  "data": [ /* array of entities */ ]
}
```

## HTTP Status Codes

| Code | Meaning | Example |
|------|---------|---------|
| 200 | OK | GET, PUT, PATCH with success |
| 201 | Created | POST with success |
| 400 | Bad Request | Invalid request body |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Resource already exists |
| 500 | Internal Server Error | Database error |

## Testing All Endpoints

### Quick Test Script

```bash
#!/bin/bash

BASE_URL="http://localhost:3000/api/v1"

# 0. Login to get token
LOGIN=$(curl -s -X POST http://localhost:3000/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}')
TOKEN=$(echo $LOGIN | jq -r '.token')
echo "✓ Logged in: $TOKEN"

# 1. Create Board
BOARD=$(curl -s -X POST $BASE_URL/boards \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Board"}')
BOARD_ID=$(echo $BOARD | jq '.data.id')
echo "✓ Created Board: $BOARD_ID"

# 2. Create List
LIST=$(curl -s -X POST $BASE_URL/boards/$BOARD_ID/lists \
  -H "Content-Type: application/json" \
  -d '{"title":"To Do"}')
LIST_ID=$(echo $LIST | jq '.data.id')
echo "✓ Created List: $LIST_ID"

# 3. Create Task
TASK=$(curl -s -X POST $BASE_URL/lists/$LIST_ID/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Sample Task"}')
TASK_ID=$(echo $TASK | jq '.data.id')
echo "✓ Created Task: $TASK_ID"

# 4. Get Everything
echo "\n✓ Getting all boards:"
curl -s -X GET $BASE_URL/boards | jq '.data | length'

echo "✓ Getting board lists:"
curl -s -X GET $BASE_URL/boards/$BOARD_ID/lists | jq '.data | length'

echo "✓ Getting list tasks:"
curl -s -X GET $BASE_URL/lists/$LIST_ID/tasks | jq '.data | length'

echo "\n✅ All endpoints working!"
```

Save as `test-api.sh` and run:
```bash
chmod +x test-api.sh
./test-api.sh
```

## Database Schema Diagram

```
users
├── id (PK)
├── username
├── email
├── password_hash
├── created_at
└── updated_at
    │
    └─── (1) ─┐
           boards(user_id → users.id)
           ├── id (PK)
           ├── user_id (FK)
           ├── title
           ├── slug
           ├── created_at
           └── updated_at
               │
               └─── (1) ─┐
                      lists(board_id → boards.id)
                      ├── id (PK)
                      ├── board_id (FK)
                      ├── title
                      ├── position
                      ├── created_at
                      └── updated_at
                          │
                          └─── (1) ─┐
                                 tasks(list_id → lists.id)
                                 ├── id (PK)
                                 ├── list_id (FK)
                                 ├── title
                                 ├── description
                                 ├── position
                                 ├── is_completed
                                 ├── due_date
                                 ├── created_at
                                 └── updated_at
```

## Load Testing Example

```bash
# Install Apache Bench
apt-get install apache2-utils  # Ubuntu/Debian
brew install httpd             # macOS

# Test endpoint concurrency
ab -n 1000 -c 10 http://localhost:3000/api/v1/boards

# Output interpretation:
# Requests per second - throughput
# Time per request - response time
# Failed requests - errors
```

## Common API Combinations

### Create a Complete Board Setup

```bash
# 1. Create board
B=$(curl -s -X POST http://localhost:3000/api/v1/boards \
  -H "Content-Type: application/json" \
  -d '{"title":"Project Roadmap"}' | jq '.data.id')

# 2. Create lists
L1=$(curl -s -X POST http://localhost:3000/api/v1/boards/$B/lists \
  -H "Content-Type: application/json" \
  -d '{"title":"Backlog"}' | jq '.data.id')

L2=$(curl -s -X POST http://localhost:3000/api/v1/boards/$B/lists \
  -H "Content-Type: application/json" \
  -d '{"title":"In Progress"}' | jq '.data.id')

L3=$(curl -s -X POST http://localhost:3000/api/v1/boards/$B/lists \
  -H "Content-Type: application/json" \
  -d '{"title":"Done"}' | jq '.data.id')

# 3. Create tasks in lists
curl -X POST http://localhost:3000/api/v1/lists/$L1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Feature 1"}'

curl -X POST http://localhost:3000/api/v1/lists/$L2/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Bug Fix 1"}'

# Done!
echo "Board $B is ready with 3 lists!"
```
