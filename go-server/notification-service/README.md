# Go Notification Service

Run locally:

```bash
cd go-server/notification-service/cmd/notification-service
go run main.go
```

Environment variables:

- `PORT`
- `API_KEY` - when set, the service requires the `X-API-Key` header for protected routes
- `SMTP_HOST`
- `SMTP_PORT`
- `SMTP_USERNAME`
- `SMTP_PASSWORD`
- `SMTP_SENDER`
- `SLACK_WEBHOOK_URL`
- `TEAMS_WEBHOOK_URL`

Health check:

```bash
curl http://localhost:8080/health
```

POST `/notify` with JSON:

```json
{
  "agent": {"name":"agent1","owner":"dev","group":"ops"},
  "title":"Alert",
  "message":"Hello",
  "emailRecipients":["user@example.com"],
  "slackWebhookUrls":["https://hooks.slack.com/..."],
  "teamsWebhookUrls":["https://outlook.office.com/webhook/..."]
}
```
