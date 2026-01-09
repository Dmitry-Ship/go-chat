# AGENTS.md

## Project Overview
Real-time multi-room chat: Go 1.23+ backend (DDD), Next.js 16 frontend, PostgreSQL/Redis/WebSocket.

## Commands

### Quick Start
```bash
make help                                  # Show all commands
make backend_build frontend_build         # Build both
make backend_test frontend_test            # Run all tests
make all_lint                              # Lint everything
make docker_up_detached                    # Start all services
```

### Backend (Go)
```bash
cd backend
go test ./...                              # All tests
go test ./internal/services -run TestLogin -v # Single test verbose
go test ./internal/... -cover             # With coverage
go test -count=1 ./...                    # Disable cache
go build -o go-bin ./cmd/server           # Build
go run ./cmd/server                        # Run server
golangci-lint run                          # Lint
golangci-lint run ./internal/... --fix     # Auto-fix
```

### Frontend (Next.js)
```bash
cd frontend
npm run dev                                # Dev server
npm run build                              # Production build
npm run lint                               # Lint
npx tsc --noEmit                           # Type check
npm test                                   # Run tests
```

## Backend (Go) Style

**Structure**: `cmd/server` → `internal/domain` → `services` → `infra` → `server`

**DDD Layers**: Domain (value objects, aggregates, repository interfaces), Infrastructure (sqlc-generated, repository impls), Application (services), Interfaces (handlers, WebSocket)

**Naming**: Files `snake_case.go`, packages lowercase. Exported: `UserRepository`, `GetByID`. Private: `userService`, `userID`. Value objects: `userName`, `userPassword`.

**Imports**: Group stdlib, third-party, local. Blank imports for side effects. Use full module paths: `"GitHub/go-chat/backend/internal/domain"`

**Formatting**: `gofmt`, tabs, error-first: `result, err := someFunc()`. No unnecessary comments.

**Context**: Always pass `context.Context` as first parameter to I/O functions. Retrieve from request: `userID := r.Context().Value(userIDKey).(uuid.UUID)`

**Error Handling**: Always check errors, early returns. Domain errors in domain layer. Custom errors: `var ErrName = errors.New("description")`. Wrap: `fmt.Errorf("create user: %w", err)`

**Value Objects**: Private struct, `NewTypeName(name string) (typeName, error)` validates and returns value object.

**Aggregates**: Root embeds `aggregate` struct. Domain events via `AddEvent()`, retrieve via `GetEvents()`.

**Repositories**: Interfaces in domain/, implementations in infra/postgres/. Methods: `Store`, `Update`, `GetByID`, `FindByUsername`, `Delete`.

**Services**: Interface with methods, private struct with repos, `NewService` constructor. Methods contain business logic only.

**Testing**: `testify/assert`. Mock structs implement repository interfaces. Test names `Test<FunctionName>`. Use `uuid.New()` for test IDs. Table-driven: define `testCase` struct, loop with `t.Run()`.

**sqlc**: Type-safe SQL in `internal/infra/postgres/db/`. Queries in `queries.sql` use `-- name: FunctionName :exec|:one|:many`.

**Database**: PostgreSQL pgx/v5 pool, Redis caching, Goose migrations in `migrations/`.

## Frontend (Next.js) Style

**Structure**: `app/` (Next.js App Router), `components/` (UI components), `hooks/` (React hooks), `contexts/` (React contexts), `lib/` (utilities)

**Naming**: Components PascalCase: `MessageItem.tsx`, `useAuth.ts`. Hooks: `use*` prefix. Utilities camelCase.

**Formatting**: ESLint with `eslint-config-next`. TypeScript strict mode enabled. Tailwind CSS for styling.

**Components**: "use client" directive for interactive components. Props interfaces inline or above component.

**React Patterns**: Context providers for global state (AuthContext, ChatContext). Custom hooks (`useAuth`, `useConversation`) encapsulate TanStack Query logic.

**API**: REST calls via `@/lib/api`, WebSocket via `@/lib/websocket`. TanStack Query for data fetching/mutations.

**Imports**: Absolute imports with `@/` alias: `@/components/ui/button`, `@/hooks/useAuth`.

**Type Safety**: Interfaces for DTOs in `@/lib/types`. Explicit return types in API functions.

**UI Components**: shadcn/ui pattern with Base UI primitives, Tailwind, class-variance-authority for variants.

## Caching

**Redis Cache**: Decorator pattern in `internal/infra/cache/`. TTL: Users 15min, Conversations 15min, Participants 10min, UserConvList 5min.

## Rate Limiting

**Sliding Window**: In-memory in `internal/ratelimit/`. Dual limits: IP-based + user-based. Returns `429 Too Many Requests` with `Retry-After` header.

## WebSocket

**Handler**: `internal/server/wsHandlers.go`. Client struct in `internal/websocket/`. Hub pattern for room broadcasts. Handle: `join`, `send_message`, `leave` events.

**Message Types**: `TextMessage`, `JoinedMessage`, `LeftMessage`, `ErrorMessage`. JSON with `type` field.

## Workflow

**Backend**: Domain → Interfaces → Implementations → Services → Wire in main.go → Test

**Pre-commit**: Run `make all_lint`, `make backend_test` before committing.

## Environment

Copy `.env.example` to `.env`. Requires DB, Redis, and JWT secrets. Run `make secret` to generate secrets.

## Common Tasks

**Add new domain model**: 1) Create value objects in `internal/domain/`, 2) Define repository interface in domain/, 3) Implement in `internal/infra/postgres/`, 4) Wrap with cache decorator in `internal/infra/cache/`, 5) Create service in `internal/services/`, 6) Wire in `cmd/server/main.go`, 7) Add sqlc queries, 8) Write tests with mocks

**Add new API endpoint**: Add handler in `internal/server/handlers.go`, register route in `cmd/server/main.go`

**Add new UI component**: Use shadcn pattern in `components/ui/`, or domain-specific in `components/chat/`

## Important Notes

- Never commit secrets or keys
- Always run `make all_lint` before committing
- Use full module paths for Go imports
- All client components need "use client" directive
