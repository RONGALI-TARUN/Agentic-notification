package main

import (
	"log"

	"github.com/example/notification-service/pkg/config"
	"github.com/example/notification-service/pkg/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	svc := server.NewServer(cfg)
	if err := svc.Run(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
