# Multi-Service Notification Repos

This workspace contains two notification service implementations:

- `go-server/notification-service` — generic Go notification service with Docker and Kubernetes support
- `java-server/production` — Maven-based Java notification service with Docker and Kubernetes support

## Structure

- `go-server/notification-service/cmd/notification-service` — Go application entrypoint
- `go-server/notification-service/pkg` — Go config, notification, server, and utility packages
- `go-server/notification-service/Dockerfile` — Go Docker build
- `go-server/notification-service/docker-compose.yml` — local compose setup
- `go-server/notification-service/k8s` — Kubernetes manifests
- `go-server/notification-service/.github/workflows` — Go CI workflow

- `java-server/production/src/main/java` — Java source tree in packages
- `java-server/production/pom.xml` — Maven build file
- `java-server/production/Dockerfile` — Java Docker build
- `java-server/production/docker-compose.yml` — local Java compose setup
- `java-server/production/k8s` — Kubernetes manifests
- `java-server/production/.github/workflows/build.yml` — Java CI workflow

## Running

### Go

```bash
cd go-server/notification-service/cmd/notification-service
go run main.go
```

### Java

```bash
cd java-server/production
mvn -B package
java -cp target/notification-service-0.1.0.jar com.example.notification.server.Application
```

## Docker

```bash
docker build -t notification-service:go ./go-server/notification-service
docker build -t notification-service-java ./java-server/production
```

## Kubernetes

```bash
kubectl apply -f go-server/notification-service/k8s
kubectl apply -f java-server/production/k8s
```
