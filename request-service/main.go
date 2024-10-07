package main

import (
	"context"
	"fmt"
	pb "go-compiler/request-service/generated/go-compiler/generated/requestpb"
	"go-compiler/request-service/internal/ports/factory"
	"go-compiler/request-service/pkg/router"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type server struct {
	pb.UnimplementedRequestServiceServer
	portFactory *factory.PortFactory
}

func NewGrpcServer() *server {
	return &server{
		portFactory: factory.NewPortFactory(),
	}
}

// Implement the Ping gRPC method
func (s *server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "pong"}, nil
}

// Implement the SubmitRequest gRPC method
func (s *server) SubmitRequest(ctx context.Context, req *pb.SubmissionRequest) (*pb.SubmissionResponse, error) {
	// You can use your portFactory's RequestController to handle the request here
	return s.portFactory.RequestController.SubmitRequest(ctx, req)

}

func main() {
	// HTTP server setup
	appRouter := router.GetRouter()
	httpPort := ":8080"

	httpServer := &http.Server{
		Addr:    httpPort,
		Handler: appRouter,
	}

	// gRPC server setup
	grpcPort := ":50051"
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		panic(fmt.Sprintf("failed to listen on port %s: %v", grpcPort, err))
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRequestServiceServer(grpcServer, NewGrpcServer())
	reflection.Register(grpcServer)

	// Run HTTP and gRPC servers concurrently
	go func() {
		fmt.Println("Starting HTTP server on port", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("HTTP server failed: %v", err))
		}
	}()

	go func() {
		fmt.Println("Starting gRPC server on port", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			panic(fmt.Sprintf("gRPC server failed: %v", err))
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("Shutting down servers...")

	// Shut down HTTP server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		fmt.Println("HTTP server forced to shutdown:", err)
	}

	// Stop gRPC server
	grpcServer.GracefulStop()

	fmt.Println("Servers stopped gracefully.")
}
