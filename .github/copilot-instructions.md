# Project Guidelines

## Code Style
- Use Go idioms and keep package APIs small and explicit.
- Prefer existing project dependencies and patterns over introducing new libraries.
- Follow existing naming and struct-tag style (for example YAML and JSON tags in config and message models).
- Keep comments concise; preserve the current mixed Chinese/English style used by the project when editing nearby code.

## Architecture
- Entry points are centered on `Bot` in [bot.go](bot.go), which composes:
  - WebSocket transport client in [ws/ws.go](ws/ws.go)
  - Event dispatcher and command center in [event/dispatcher.go](event/dispatcher.go) and [event/command.go](event/command.go)
  - OneBot API sender and response parsing in [api/sender.go](api/sender.go) and [api/resp.go](api/resp.go)
- Message representation uses `message.Chain` and `message.Segment` with CQ parsing/escaping in [message/cq.go](message/cq.go).
- NapCat-specific API extensions live in [extensions/napcat/napcat.go](extensions/napcat/napcat.go).

## Build And Test
- Install dependencies: `go mod download`
- Run all tests: `go test ./...`
- Run a focused package test: `go test ./ws -run TestClient`

## Conventions
- `Bot.Start()` is non-blocking. When adding examples or runtime flows, ensure the process is explicitly blocked and shutdown is handled like [examples/simple/main.go](examples/simple/main.go).
- WebSocket reconnect semantics are important:
  - `Ws.MaxReconnect` default is `3`
  - `0` disables automatic reconnect
  - Reconnect supervision and stop logic are implemented in [ws/ws.go](ws/ws.go)
- Event handling order is significant: command processing first (message events), then all-type handlers, then type-specific handlers. Preserve this behavior when changing dispatch flow.
- Be careful with command propagation logic in [event/command.go](event/command.go): only call interfaces that are actually implemented.
- CQ string escaping/unescaping must stay compatible with existing tests in [message/cq_test.go](message/cq_test.go) and command parsing tests in [event/command_test.go](event/command_test.go).

## Testing Expectations
- Add or update tests for behavior changes in the touched package.
- Prefer `testify` patterns already used across the repo (`assert`/`require`).
- For transport or retry behavior, mirror patterns from [ws/ws_test.go](ws/ws_test.go) using `httptest` and `require.Eventually`.