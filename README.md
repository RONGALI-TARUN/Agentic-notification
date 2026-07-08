# Attention Orchestration Platform

This workspace now contains the initial milestone-1 foundation for an attention orchestration platform.

## Included in milestone 1

- A Go API service with a health endpoint at /health
- A lightweight React + TypeScript web shell for the future UI
- Local PostgreSQL via Docker Compose
- CI wiring for the Go service

## Structure

- [go-server/notification-service](go-server/notification-service) — Go API service and Docker setup
- [apps/web](apps/web) — Vite + React + TypeScript frontend shell
- [docs/milestone-1.md](docs/milestone-1.md) — milestone plan and review notes

## Run locally

```bash
docker compose up --build
```

Then open:

- API: http://localhost:8080/health
- UI: http://localhost:3000

