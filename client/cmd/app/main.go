package main

import (
	"github.com/mrkucher83/hash-service/client/internal/routes"
	"github.com/mrkucher83/hash-service/client/pkg/logger"
	"syscall"
)

const DefaultPort = "8080"

func main() {
	logger.InitLogger(logger.NewLogrusLogger())

	port := DefaultPort
	if value, ok := syscall.Getenv("HASHER_PORT"); ok {
		port = value
	}

	routes.Start(port)
}
