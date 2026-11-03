# Go scaffold (AUT-6)

Authopus is a **single Go module** (see `go.mod`) and a **modular monolith** (ADR-001). The two command entry points share internal packages; they are not separate services. This milestone intentionally contains no real authorization, identity, audit, database, or MCP functionality.

## Requirements

- Go 1.23 or newer compatible Go toolchain
- Make for convenience targets
- Docker for an optional image build

## Commands

```sh
go run ./cmd/authopus       # Hello World!
go run ./cmd/authopusctl    # Hello World!
make build                  # builds bin/authopus and bin/authopusctl
make fmt-check              # validates gofmt formatting
make vet                    # static analysis
make test                   # go test ./...
docker build -t authopus:dev .
docker run --rm authopus:dev
```

Configuration: **none yet**. Both commands are deterministic Hello World placeholders and do not open listeners, connect to PostgreSQL, or read secrets. Do not deploy this milestone as a functioning security gateway. CI runs formatting, vet, tests, build, and executable smoke checks on pushes and pull requests.

## Planned package boundaries

- `internal/enforcement` — Gate: MCP enforcement point
- `internal/identity` — Delegate: trusted identity and delegation
- `internal/authorization` — Guard: authorization orchestration
- `internal/policy` — Guard policy evaluation and Codex policy lifecycle (to be separated by interfaces as features arrive)
- `internal/auditapproval` — Watchtower: audit/approval workflows
- `internal/platform/{config,database,httpserver,observability}` — shared infrastructure
- `internal/app` — composition and lifecycle

The remaining directories reserve locations for contracts, example policies, migrations, and tests. Their presence does **not** imply implemented functionality. Avoid inventing APIs/policies/migrations before their respective design issues are completed.
