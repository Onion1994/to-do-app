# Todo App

Simple todo app in Go. Supports CLI and HTTP server. Concurrent-safe file storage.

## Features

- CLI and HTTP server  
- CRUD todos  
- Concurrent-safe file operations (Actor/CSP pattern)  
- Graceful shutdown  
- Basic logging  
- Unit and concurrency tests  

## Installation

```bash
git clone https://github.com/Onion1994/to-do-app.git
cd todo-app
go mod download
go build
```

## Usage

### CLI Mode
```bash
./todo-app -view
./todo-app -add "buy milk"
./todo-app -remove "buy milk"
./todo-app -find "buy groceries" -update-status "completed"
./todo-app -find "buy groceries" -update-description "buy milk"
```

### Server Mode

Start the HTTP server:
```bash
./todo-app -mode server
```

Server runs on `http://localhost:8080`

Press `Ctrl+C` to gracefully shutdown the server.

## API Endpoints

### Create Todo
```http
POST /create
Content-Type: application/json

{
  "description": "buy groceries"
}
```

### Read All Todos
```http
GET /read
```

### Update Todo
```http
PATCH /update
Content-Type: application/json

{
  "description": "buy groceries",
  "field": "status",
  "newValue": "completed"
}
```

### Delete Todo
```http
DELETE /delete
Content-Type: application/json

{
  "description": "buy groceries"
}
```

### Web Interface
```http
GET /list        # View todos in browser
GET /about/      # Static about page
```

## Status Values

Valid todo statuses:
- `not started` (default)
- `started`
- `completed`

## Testing

### Run All Tests
```bash
go test ./...
```