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
```

### Backend (Go)
```bash
cd backend

# Tests
go test ./...                              # All tests
go test ./internal/services -run TestLogin # Single test
go test ./internal/... -v                  # Verbose for package
go test -count=1 ./...                     # Disable cache

# Build & Lint
go build -o go-bin ./cmd/server
golangci-lint run
```

### Frontend (TypeScript/Next.js)
```bash
cd frontend

npm run dev                               # Dev server
npm run build                             # Production build
npm run start                             # Production server
npm run lint                              # ESLint
npx tsc --noEmit                          # Type check
```

## Backend (Go) Style

**Structure**: `cmd/server` → `internal/domain` (models/interfaces) → `services` → `infra` (postgres/redis) → `server` (http/ws)

**DDD Layers**: Domain (value objects, aggregates, repository interfaces), Infrastructure (GORM models, repository impls), Application (services), Interfaces (handlers, WebSocket)

**Naming**: Files `snake_case.go`, packages lowercase. Exported: `UserRepository`, `GetByID`. Private: `userService`, `userID`. Value objects: `userName`, `userPassword`.

**Imports**: Group stdlib, third-party, local. Blank imports for side effects. Use full module paths: `"GitHub/go-chat/backend/internal/domain"`

**Formatting**: `gofmt`, tabs, error-first: `result, err := someFunc()`. No unnecessary comments.

**Value Objects**: Private struct, `NewTypeName(name string) (typeName, error)` validates, returns value object. Use for primitives like username, password.

**Aggregates**: Root embeds `aggregate` struct. Methods modify aggregate state.

**Repositories**: Interfaces in domain/, implementations in infra/postgres/. Methods: `Store`, `Update`, `GetByID`.

**Services**: Interface with methods, private struct with repos, `NewService` constructor. Methods contain business logic only.

**Testing**: `testify/assert`. Mock structs implement repository interfaces. Test names `Test<FunctionName>`. Use `uuid.New()` for test IDs. Mock `methodsCalled` map for verification.

**GORM**: Models in `internal/infra/postgres/models.go`, tags `gorm:"primaryKey;type:uuid"`. Use repository pattern, never GORM directly in services.

**HTTP/WebSocket**: Handlers in `internal/server/`. Handlers accept services, wire in `main.go`. WebSocket in `wsHandlers.go`, client management in `internal/websocket/`.

**Error Handling**: Always check errors, early returns. Domain errors in domain layer (e.g., "username is empty"). Descriptive messages.

## Frontend (TypeScript/Next.js) Style

**Structure**: `app/` (App Router with `(auth)`/`(authenticated)` groups), `src/components/`, `api/`, `contexts/`, `types/`

**Naming**: Files `PascalCase.tsx` (components), `camelCase.ts` (utils). Components PascalCase, functions/vars camelCase, types PascalCase, constants UPPER_SNAKE_CASE.

**Imports**: React hooks → third-party → local. Use `"react"` hooks, relative paths for local: `"../types/coreTypes"`.

**TypeScript**: Strict mode. Props: `type Props = { user: User }`. Generics: `makeQuery<T>(url)`. No `any`. Type assertions: `as T`.

**Components**: Functional with hooks. `"use client"` for interactivity. Destructure props. Link from next/link.

**Styling**: CSS Modules: `import styles from "./Component.module.css"`. Class names: `styles.wrap`, `styles.conversationInfo`.

**API**: Centralized in `src/api/fetch.ts`. Factory functions: `makeCommand<T>(url)` returns `async (body?: T) => any`, `makeQuery<T>(url)` returns `(param) => () => Promise<T>`. Axios with `/api` base URL.

**State**: React Context for auth (global), React Query for server state, `useState`/`useEffect` for component state, WebSocket context.

**Error Handling**: Try-catch async operations. Display via `ErrorAlert` component. `isChecking` boolean for loading states.

## Caching

**Redis Cache**: Centralized in `internal/infra/cache/`. Decorator pattern wraps repositories: `UserCacheDecorator`, `GroupConversationCacheDecorator`, `ParticipantCacheDecorator`.

**Cache Keys**: Defined in `cacheKeys.go` with prefixes (user, conv, participants, conv_meta, user_conv_list). Keys: `UserKey(id)`, `UsernameKey(name)`, `ParticipantsKey(conversationID)`.

**TTL**: Users (15min), Conversations (15min), Participants (10min), ConvMeta (15min), UserConvList (5min).

**Invalidation**: Event-driven via `CacheInvalidationService`. Subscribes to domain events (renamed, deleted, joined, left, invited). Invalidates by key pattern on mutations.

**Usage**: Wire in `cmd/server/main.go`. Use `cached*Repository` instead of raw repos. Decorators implement repository interfaces transparently.

**Serialization**: JSON serialization for domain models (`SerializeUser`, `DeserializeUser`). Passwords replaced with dummy value in cache.

## Rate Limiting

**Sliding Window**: In-memory implementation in `internal/ratelimit/`. `Config` struct: `MaxConnections`, `WindowDuration`. `NewSlidingWindowRateLimiter(config)` creates limiter.

**Middleware**: `wsRateLimit` in `internal/server/rateLimitMiddleware.go`. Applies to WebSocket connections. Dual limits: IP-based + user-based (if authenticated).

**IP Detection**: Extracts from `X-Forwarded-For`, `X-Real-IP`, or `RemoteAddr`. Supports proxy headers for proper IP resolution.

**Response**: Returns `429 Too Many Requests` with `Retry-After` header (seconds until next slot).

**Testing**: Mock rate limiters with specific config. Tests verify within limit, exceeded limit, sliding window expiration, header parsing.

## Workflow

**Backend changes**: Domain → Interfaces in domain/ → Implementations in infra/ → Services → Wire in cmd/server/main.go → Test with mocks

**Frontend changes**: Types in types/ → API in fetch.ts → Context if global → Component in app/ → Test (manual, no test suite)

**Commit**: Run `make all_lint`, `make backend_test` before committing. Fix all failures.

## Environment

Copy `.env.example` to `.env`. Backend: DB_HOST, DB_PORT, DB_USER, DB_NAME, DB_PASSWORD, REDIS_HOST, REDIS_PORT, REDIS_PASSWORD, ACCESS_TOKEN_SECRET, REFRESH_TOKEN_SECRET, CLIENT_ORIGIN. Frontend proxies to `/api`.
