# Task API

Small Go HTTP API for managing in-memory tasks. Data is lost when the server
stops. Supports GET, POST, PATCH, DELETE requests.

## Run

Requires Go 1.27+.

```bash
go run ./cmd/server
```

The server listens on `http://localhost:8080`.

## Tests

Run:

```bash
go test ./...
```

This discovers and runs all Go tests in every package under the project.
`internal/api/handlers_test.go` contains HTTP handler unit tests using
`httptest`, covering task creation, title validation, retrieval, completion
updates, deletion, and not-found responses.

## Endpoints

| Method | Path | Description |
|---|---|---|
| GET | `/health` | Health check |
| GET | `/tasks` | List tasks; optional case-insensitive `?q=` title filter |
| POST | `/tasks` | Create a task |
| GET | `/tasks/{id}` | Get a task |
| PATCH | `/tasks/{id}` | Update completion status |
| DELETE | `/tasks/{id}` | Delete a task |

Create a task:

```bash
curl -X POST http://localhost:8080/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title":"Try the API"}'
```

Task format:

```json
{"id":1,"title":"Try the API","done":false}
```

`title` is required, trimmed, and limited to 144 Unicode characters.
Invalid input returns `400`; missing tasks return `404`; successful deletion
returns `204`.

## Structure

```text
cmd/server/main.go          Server and route setup
internal/api/handlers.go   HTTP handlers
internal/api/middleware.go Logging and CORS
internal/api/responses.go  JSON response helpers
internal/storage/memory.go Thread-safe in-memory store
```
