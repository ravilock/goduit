# Agent Guidelines for goduit

This repository is a Go implementation of Conduit (a simple social reading API). It uses Docker Compose for running services with RabbitMQ or Redis queue backends. Follow these guidelines when working on this codebase.

## Quick Start

```bash
# Run unit tests
go test $(go list ./... | grep -v integrationTests) -count=1

# Run lint
./bin/golangci-lint run --timeout 5m

# Setup Golang CI Lint (first time only)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s v2.1.6

# Run integration tests (requires Docker Compose services running)
make test-integration
```

### Running a Single Test

```bash
# Basic single test run
go test ./path/to/package -run TestSpecificTest

# Single test with verbose output
go test ./path/to/package -v -run TestSpecificTest

# Single test in CI-like mode (no parallel execution)
go test ./path/to/package -count=1 -run TestSpecificTest

# Full package tests without parallelization
go test $(go list ./path/to/package/...)/... -count=1
```

## Build & Run Commands

| Command | Description |
|---------|-------------|
| `make run` | Start all services (RabbitMQ by default) |
| `make run QUEUE=redis` | Start with Redis queue |
| `make logs-api` | View API service logs |
| `make logs-worker` | View feed worker logs |
| `make stop-all` | Stop all services |
| `make test` | Run unit tests inside container |
| `make lint-check` | Run linter locally |

### Queue Configuration

```bash
# RabbitMQ (default)
QUEUE_TYPE=rabbitmq
QUEUE_URL=amqp://guest:guest@goduit-queue:5672/

# Redis
QUEUE_TYPE=redis
QUEUE_URL=redis://goduit-redis:6379
```

## Code Style Guidelines

### Imports

**Order and Grouping:**
1. Standard library imports (blank import if needed)
2. Third-party imports (no package name)
3. Local imports (package name, sorted alphabetically)

Example:
```go
import (
    "context"
    "errors"
    "log/slog"

    "github.com/labstack/echo/v4"
    "github.com/ravilock/goduit/api"
    articleHandlers "github.com/ravilock/goduit/internal/articlePublisher/handlers"
)
```

### Formatting

- Run `gofmt -s` or rely on golangci-lint
- Each struct field and function must have a comment (Go standard)
- Blank lines between logical blocks
- No trailing whitespace
- Single space around operators and after commas

### Naming Conventions

**Types:** PascalCase (`DeleteCommentHandler`, `AppError`, `Server`)
**Functions/Methods:** camelCase (`NewDeleteCommentHandler`, `Start()`, `Error()`)
**Variables/Constants:** camelCase (`articlePublisher`, `queueName`, `QUEUE_TYPE`)
**Interface Methods:** camelCase (even in method signatures)

### Struct Naming

```go
// Use descriptive names with clear singular/plural patterns
type Server struct {                    // Single instance
    Echo   *echo.Echo       // Named field, not anonymous if used directly
    db     *mongoDriver.Client
    queue  queue.Connection
}

type DeleteCommentHandler struct {      // Action handler
    commentDeleter   commentDeleter
    commentGetter    commentGetter
    articlePublisher articleGetter
}
```

### Interfaces

**Naming:** Add domain context to prefixed interfaces
**Receivers:** Keep method names lowercase, use receiver type name as doc reference

```go
type commentDeleter interface {
    DeleteComment(ctx context.Context, ID string) error
}

func (h *DeleteCommentHandler) DeleteComment(c echo.Context) error {
    // Document how the handler uses each interface in comments
}
```

### Type Assertion Pattern

**Use `errors.As()` for AppError type assertion:**

```go
if appError := &app.AppError{}; errors.As(err, &appError) {
    switch appError.ErrorCode {
    case app.ArticleNotFoundErrorCode:
        return api.ArticleNotFound(request.Slug)
    case app.CommentNotFoundErrorCode:
        return api.CommentNotFound(request.ID)
    }
}
```

**This pattern preserves the original error for logging while handling specific cases.**

### Error Handling

**API Layer:** Create HTTPError for client-facing responses
**App Layer:** Create AppError with ErrorCode for internal logic

```go
// App layer - defines error type and code
func ArticleNotFoundError(identifier string) *AppError {
    return &AppError{
        ErrorCode:     ArticleNotFoundErrorCode,
        CustomMessage: fmt.Sprintf("Article with identifier %q was not found", identifier),
    }
}

// API layer - HTTP response mapping
func ArticleNotFound(identifier string) *echo.HTTPError {
    return &echo.HTTPError{
        Code:    http.StatusNotFound,
        Message: fmt.Sprintf("Article with identifier %q not found", identifier),
    }
}
```

**Pre-defined Errors (api/errors.go):**
- `FailedLoginAttempt`, `FailedAuthentication` — authentication failures
- `ConfictError` — content collision
- `Forbidden`, `CouldNotUnmarshalBodyError` — validation/middleware errors
- `UserNotFound(id)`, `ArticleNotFound(id)`, `CommentNotFound(id)` — resource not found

### Context Handling

```go
ctx := c.Request().Context()
serviceMethod(ctx, args...)
```

Always extract context from request before calling services. Don't pass `nil` or ignore context.

### Mocking

**Use mockery for generating interface mocks:**

```bash
mockery --all
```

**Test file patterns:**
- Interfaces: lowercase type names with domain prefix (`commentDeleter`, not `IComment`)
- Test files end with `_test.go`
- Use testify `require` assertions
- Setup validator once per test setup: `err := validators.InitValidator(); require.NoError(t, err)`

### Code Comments (Documentation)

```go
// Server represents the API server implementation. It manages Echo instance,
// database connection, and queue publisher for article feed processing.
type Server interface {
    http.Handler
    Start()
}

// writeCommentHandler handles POST /api/articles/:slug/comments requests.
// It validates that:
//   - The user owns the comment (authorization check)
//   - The associated article exists
func (h *WriteCommentHandler) WriteComment(c echo.Context) error {
```

## Testing Guidelines

### Test Structure

1. `TestFunctionName` — outer function covering scenario
2. `t.Run("Sub-test-name", func(t *testing.T) { ... })` — sub-scenarios
3. TDD pattern: Arrange → Act → Assert blocks
4. Sub-tests for different input states, edge cases, and error paths

### Test Dependencies

- Unit tests: mocked dependencies (generate with `mockery --all`)
- Integration tests: real MongoDB, queue (RabbitMQ/Redis) in Docker Compose

### Testing Checklist

- [ ] Handle success path
- [ ] Test unauthorized access
- [ ] Test resource not found
- [ ] Test validation errors
- [ ] Test edge cases (empty input, limits)
- [ ] Verify correct HTTP status codes

## Linting

```bash
# Check all linters
./bin/golangci-lint run --timeout 5m

# Fix what it can
./bin/golangci-lint run --fix --timeout 5m
```

### golangci-lint Configuration (`.golangci.yml`)

- Default linters: `standard` set
- errcheck excluded for test files and integration tests
- Warn on unused code (unless intentionally unexported)

## Architecture Overview

### Layering

```
api/                    # HTTP handlers, Echo middleware, JSON responses
internal/app/           # Domain errors, interfaces
internal/articlePublisher/ # Core business logic (services, handlers within subpackages)
  - assemblers/,   # Data transformation
    handlers/,     # Service layer for each endpoint
    models/,       # MongoDB documents
    publishers/,   # Queue publishing
    requests/,     # Validation request objects
    repositories/, # DB operations
    services/      # Business logic
internal/followerCentral/  # Similar structure to articlePublisher
internal/profileManager/
internal/identity/        # JWT auth, middleware
internal/log/             # Logger setup
internal/mongo/           # Database client wrapper
internal/queue/           # Queue abstraction (RabbitMQ/Redis agnostic)
```

### Handler Pattern

Each endpoint gets a dedicated handler struct:

```go
type WriteCommentHandler struct {
    writeCommentService     writeCommentService       // Business logic
    getProfileService       getProfileService         // For authorship check
}

func NewWriteCommentHandler(
    service writeCommentService,
    profileService getProfileService,
) *WriteCommentHandler
```

This promotes single responsibility and testability.

### Service Pattern

Services operate on repositories (no dependencies):

```go
type WriteCommentService struct {
    commentRepository articlePublisher.Repository // DB access only
}

func NewWriteCommentService(repository articlePublisher.Repository) *WriteCommentService
```

No HTTP context, no cookies—pure business logic.

## Environment Variables

Set these in Docker Compose or `.env`:

```bash
PORT=5000                          # API listen port
DB_URL=mongodb://mongo:27017/      # MongoDB connection string
QUEUE_TYPE=rabbitmq                # rabbitmq | redis
QUEUE_URL=amqp://guest:guest@...   # Queue connection string
```

JWT keys are generated at runtime via `scripts/generateJWTRS256Keys.sh` for RS256-signed tokens.

## Common Tasks

### Adding a New Endpoint

1. Create handler struct and constructor in `internal/articlePublisher/handlers/`
2. Define service interface and implementation if new logic
3. Register Echo route in API server initialization
4. Add validation to `/api/validators/common.go` if needed
5. Write unit tests with mocked dependencies
6. Test integration with both queue types

### Adding a New Interface

1. Define in `internal/package/interface_name.go` (lowcase type name)
2. Implement concrete type following interface contract
3. Register implementation in container construction (`server.go`)
4. Generate mocks if needed for testing

### Switching Queue Backend

Restart with different QUEUE_TYPE environment variable—no code changes needed. The queue package abstracts RabbitMQ vs Redis connections.

## Additional Notes

- Use `go list ./... | grep -v integrationTests` to exclude integration tests
- Integration tests run against Docker Compose services via Makefile targets
- Generated mocks follow Veokra/mockery conventions, placed in same directory as interface
- For CI CD: lint and unit tests are mandatory; integration tests matrix runs both queue types
