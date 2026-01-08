# AGENTS.md

## Project Overview
Real-time multi-room chat: Go 1.18+ backend (DDD), Next.js 13 frontend, PostgreSQL/Redis/WebSocket.

## Commands

### Quick Start (Makefile)
```bash
make help                                  # Show all commands
make backend_build frontend_build         # Build both
make backend_test frontend_test            # Run all tests
make all_lint                              # Lint everything
make docker_up_detached                    # Start all services
make clean                                 # Clean build artifacts
```

### Backend (Go)
```bash
cd backend

# Tests
go test ./...                              # All tests
go test ./internal/services -run TestLogin -v # Single test verbose
go test ./internal/... -cover             # With coverage
go test -count=1 ./...                    # Disable cache
go test ./internal/services -run TestServiceFunctionName # Exact test

# Table-driven tests
go test ./internal/services -run TestCreateUser -v # Runs all table cases

# Build & Lint
go build -o go-bin ./cmd/server
go run ./cmd/server                        # Run server directly
golangci-lint run
golangci-lint run ./internal/... --fix     # Auto-fix issues
```

### Frontend (TypeScript/Next.js)
```bash
cd frontend

npm run dev                               # Dev server
npm run build                             # Production build
npm run start                             # Production server
npm run lint                              # ESLint
npm run lint -- --fix                     # Auto-fix ESLint
npx tsc --noEmit                          # Type check
```

## Backend (Go) Style

**Structure**: `cmd/server` → `internal/domain` (models/interfaces) → `services` → `infra` (postgres/redis) → `server` (http/ws)

**DDD Layers**: Domain (value objects, aggregates, repository interfaces), Infrastructure (sqlc-generated models, repository impls), Application (services), Interfaces (handlers, WebSocket)

**Naming**: Files `snake_case.go`, packages lowercase. Exported: `UserRepository`, `GetByID`. Private: `userService`, `userID`. Value objects: `userName`, `userPassword`.

**Imports**: Group stdlib, third-party, local. Blank imports for side effects. Use full module paths: `"GitHub/go-chat/backend/internal/domain"`

**Formatting**: `gofmt`, tabs, error-first: `result, err := someFunc()`. No unnecessary comments.

**Context**: Always pass `context.Context` as first parameter to functions that do I/O. Use for cancellation, timeouts, request-scoped values. Retrieve from request context: `userID := r.Context().Value(userIDKey).(uuid.UUID)`

**Error Handling**: Always check errors, early returns. Domain errors in domain layer (e.g., "username is empty"). Custom errors: `var ErrName = errors.New("description")`. Return wrapped errors: `fmt.Errorf("create user: %w", err)`.

**Value Objects**: Private struct, `NewTypeName(name string) (typeName, error)` validates, returns value object. Use for primitives like username, password.

**Aggregates**: Root embeds `aggregate` struct. Methods modify aggregate state. Domain events via `AddEvent()`, retrieve via `GetEvents()`.

**Repositories**: Interfaces in domain/, implementations in infra/postgres/. Methods: `Store`, `Update`, `GetByID`, `GetByUsername`, `Delete`.

**Services**: Interface with methods, private struct with repos, `NewService` constructor. Methods contain business logic only.

**Testing**: `testify/assert`. Mock structs implement repository interfaces. Test names `Test<FunctionName>`. Use `uuid.New()` for test IDs. Table-driven tests: define `testCase` struct, loop with `t.Run()`.

**sqlc**: Type-safe SQL code generation in `internal/infra/postgres/db/`. Queries in `queries.sql` use `-- name: FunctionName :exec|:one|:many` comments.

**Database**: PostgreSQL with pgx/v5 pool. Redis for caching. Migrations via Goose in `migrations/`.

## Frontend (TypeScript/Next.js) Style

**Structure**: `app/` (App Router with `(auth)`/`(authenticated)` groups), `src/components/`, `api/`, `contexts/`, `types/`

**Naming**: Files `PascalCase.tsx` (components), `camelCase.ts` (utils). Components PascalCase, functions/vars camelCase, types PascalCase, constants UPPER_SNAKE_CASE.

**Imports**: React hooks → third-party → local.

**TypeScript**: Strict mode. Props types, generics, no `any`.

**Components**: Functional with hooks. `"use client"` for interactivity. Destructure props.

**API**: Centralized in `src/api/fetch.ts`. Factory functions for commands and queries.

**State**: React Context for auth, React Query for server state, useState/useEffect for local. Query keys: `["resource", id]` format.

## Caching

**Redis Cache**: Decorator pattern in `internal/infra/cache/`. TTL: Users 15min, Conversations 15min, Participants 10min, UserConvList 5min.

## Rate Limiting

**Sliding Window**: In-memory implementation in `internal/ratelimit/`. Dual limits: IP-based + user-based. Returns `429 Too Many Requests` with `Retry-After` header.

## Workflow

**Backend**: Domain → Interfaces → Implementations → Services → Wire in main.go → Test

**Frontend**: Types → API → Context → Component

**Pre-commit**: Run `make all_lint`, `make backend_test` before committing.

## WebSocket

**Handler**: `internal/server/wsHandlers.go` manages connections. Client struct in `internal/websocket/`. Broadcast messages to rooms via hub pattern. Handle: `join`, `send_message`, `leave` events.

**Message Types**: `TextMessage`, `JoinedMessage`, `LeftMessage`, `ErrorMessage`. Serialized as JSON with `type` field for discriminated union.

## Environment

Copy `.env.example` to `.env`. Requires DB, Redis, and JWT secrets.

## Common Tasks

**Add new domain model**:
1. Create value objects in `internal/domain/`
2. Define repository interface in `internal/domain/`
3. Implement in `internal/infra/postgres/`
4. Wrap with cache decorator in `internal/infra/cache/`
5. Create service in `internal/services/`
6. Wire in `cmd/server/main.go`
7. Add sqlc queries in `internal/infra/postgres/queries.sql`
8. Write tests with mocks

**Add new API endpoint**:
1. Add handler in `internal/server/handlers.go`
2. Register route in `cmd/server/main.go`
3. Add frontend API function in `src/api/fetch.ts`
4. Create React Query hook in `src/hooks/`
5. Update frontend types in `src/types/`

## Important Notes

- Never commit secrets or keys to the repository
- Always run `make all_lint` before committing
- Use full module paths for Go imports, not relative paths
