# Speakeasy Backend — Development Guidelines

## Commands
- Run: `go run main.go`
- Build: `go build ./...`
- Test: `go test ./...`
- Vet: `go vet ./...`
- Format: `gofmt -l .`
- Tidy dependencies: `go mod tidy`

## Code Style Guidelines
- **Architecture:** Go REST API in layers — `internal/handler` → `internal/service` → `internal/repository` → `internal/model`, wired together in `internal/router`.
- **Git Commits:** Emoji + British English description (e.g. `🐛 Fix audio state machine transition`).

## Related Repositories
- Frontend: [speakeasy-frontend](https://github.com/Edgarmontenegro123/speakeasy-frontend)
