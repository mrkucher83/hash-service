package main

import (
	"github.com/mrkucher83/hash-service/client/pkg/logger"
	"github.com/mrkucher83/hash-service/server/internal/handlers"
	"github.com/mrkucher83/hash-service/server/pkg/pb"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const DefaultPort = "50051"

func main() {
	logger.InitLogger(logger.NewLogrusLogger())

	port := DefaultPort
	if value, ok := syscall.Getenv("SERVER_PORT"); ok {
		port = value
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		logger.Fatal("failed to listen a server: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStringHashServiceServer(grpcServer, &handlers.Server{})

	go func() {
		logger.Info("grpc-server launched on port: %s", port)
		if err = grpcServer.Serve(listener); err != nil {
			logger.Fatal("failed to launch a server: %v", err)
		}
	}()

	<-stop
	grpcServer.GracefulStop()
	logger.Info("gRPC server gracefully stopped")
}
