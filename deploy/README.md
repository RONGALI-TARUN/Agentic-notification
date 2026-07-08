# Local deployment

Use these local services for the current milestone:

- PostgreSQL for event persistence
- Ollama for local model-backed learning summaries

Run:

```bash
docker compose -f deploy/docker-compose.local.yml up -d
```

Then start the backend service and optionally the web app.
