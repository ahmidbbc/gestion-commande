# CLAUDE.md — project conventions

This project was scaffolded from scratch with a **clean / hexagonal** architecture
imposed by default. Keep the three layers strictly separated.

## Layers (Go)

```
Handlers (I/O)          Usecases (business)        Repositories (data)
──────────────          ───────────────────        ───────────────────
internal/api/           internal/service/          internal/store/
main.go (wiring only)
```

**Rules**
- `main.go` wires dependencies only — no business logic.
- `internal/api/` decodes input, calls a usecase, encodes output — ZERO logic.
- `internal/service/` holds all business logic; depends on repository interfaces.
- `internal/store/` is data access only — no business logic.
- Dependencies flow inward: nothing imports `api/`.
