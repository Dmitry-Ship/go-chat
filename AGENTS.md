# AGENTS.md

## Build, Lint, and Test Commands

### Backend (Go 1.23)
- **Build**: `make backend_build` or `cd backend && go build -o go-bin ./cmd/server`
- **Test all**: `make backend_test` or `cd backend && go test -timeout 2m ./...`
- **Test single**: `cd backend && go test -run TestFunctionName ./path/to/package`
- **Test with coverage**: `cd backend && go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html`
- **Test verbose**: `cd backend && go test -v ./...`
- **Lint**: `make backend_lint` or `cd backend && golangci-lint run`
- **Run dev**: `make backend_run` or `cd backend && go run ./cmd/server`

### Frontend (Next.js 16, React 19, TypeScript)
- **Build**: `make frontend_build` or `cd frontend && npm run build`
- **Test**: `make frontend_test` or `cd frontend && npm test`
- **Test single**: `cd frontend && npm test -- fileName.test.ts`
- **Test watch**: `cd frontend && npm test -- --watch`
- **Lint**: `make frontend_lint` or `cd frontend && npm run lint`
- **Type check**: `make frontend_type` or `cd frontend && npx tsc --noEmit`
- **Dev server**: `make frontend_dev` or `cd frontend && npm run dev`

### Combined Commands
- `make all_build` - Build backend and frontend
- `make all_test` - Run all tests
- `make all_lint` - Run all linters (backend + frontend)
- `make clean` - Remove build artifacts

### Docker
- `make docker_up` - Start all services
- `make docker_up_detached` - Start in detached mode
- `make docker_down` - Stop all services
- `make docker_build` - Build all Docker images

---

## Code Style Guidelines

### Backend (Go 1.23)

**Package Structure** (Clean Architecture):
- `internal/domain/` - Domain entities, business logic, validation (no dependencies)
- `internal/services/` - Application services (auth, conversation, notifications)
- `internal/infra/postgres/` - PostgreSQL repositories using sqlc/pgx
- `internal/infra/cache/` - Redis cache decorators (decorator pattern)
- `internal/server/` - HTTP handlers, WebSocket handlers, middleware
- `internal/websocket/` - WebSocket connection management
- `internal/ratelimit/` - Rate limiting (token bucket)
- `internal/readModel/` - Query models and DTOs
- `internal/config/` - Configuration loading
- `cmd/server/` - Application entry point

**Imports**: Three blocks - stdlib, third-party, internal - separated by blank lines

**Naming**: PascalCase for exported, camelCase for private, package names lowercase

**Error Handling**: Wrap with `fmt.Errorf("operation: %w", err)`, use `errors.New` only for sentinel errors

**Interfaces**: Define in domain packages, implement in infra packages (e.g., `UserRepository` in domain)

**Testing**: Use `github.com/stretchr/testify/assert`, test files `*_test.go`, table-driven tests with subtests (`t.Run`)

**Database**: sqlc-generated queries in `internal/infra/postgres/db/`, pgx for connection pooling

**Security**: bcrypt cost 14, JWT tokens, refresh token rotation, input sanitization with bluemonday

**Context**: Pass `context.Context` as first parameter to all service/repository methods

**Repository Pattern**: Decorate with cache layer when appropriate (e.g., `userCacheDecorator`)

**Middleware**: Chain using alice pattern (auth, rate limiting, logging)

---

### Frontend (Next.js 16, React 19, TypeScript)

**File Structure**:
- `app/` - Next.js App Router pages and layouts
- `components/` - Reusable UI (chat/, ui/)
- `contexts/` - React context providers (AuthContext, ChatContext)
- `hooks/` - Custom hooks
- `hooks/queries/` - TanStack Query hooks for data fetching
- `hooks/mutations/` - TanStack Query mutations
- `lib/` - Utilities and API client

**Imports**: Use path aliases `@/` for internal imports, no bare relative imports

**Components**:
- Client components: `"use client"` at top
- Named exports preferred: `export const ComponentName = ...`
- Props interface: `interface Props { prop: type }`
- Early returns for conditions, avoid nested ternaries

**TypeScript**: Strict mode, no `any`, use generics, proper types for API responses

**State Management**:
- TanStack Query for server state (queries/mutations)
- React Context for auth/chat state
- Local state (`useState`, `useReducer`) for simple UI

**Styling**: Tailwind CSS, classnames with `clsx` + `tailwind-merge`, shadcn/ui components

**Error Handling**: `try/catch` in async functions, display user-friendly messages, error boundaries

**API Calls**: Centralized in `lib/api.ts`, typed return values, `credentials: "include"` for cookies

**WebSocket**: Reconnection logic in `lib/websocket.ts`, message handling via React Context

**Environment**: Define in `.env.example`, `NEXT_PUBLIC_` prefix for client-side vars
