package main

import (
	"fmt"
	"go-compiler/request-service/internal/ports/factory"
	"go-compiler/request-service/pkg/router"
	pb "go-compiler/request-service/proto"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

func main() {
	portFactory := factory.NewPortFactory()

	// Start HTTP server
	go startHTTPServer(router.GetRouter())

	// Create and start gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterRequestServiceServer(grpcServer, portFactory.GRPCServer)
	startGRPCServer(grpcServer)
}

func startHTTPServer(appRouter http.Handler) {
	port := ":8080"
	httpServer := &http.Server{
		Addr:    port,
		Handler: appRouter,
	}

	fmt.Println("HTTP Server is running on port", port)
	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func startGRPCServer(grpcServer *grpc.Server) {
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}

	fmt.Println("gRPC Server is running on port", port)
	if err := grpcServer.Serve(lis); err != nil {
		panic(fmt.Sprintf("failed to serve: %v", err))
	}
}
