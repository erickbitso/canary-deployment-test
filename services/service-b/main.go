package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "service-b/service-b/hello" // Replace with the correct import path

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Server is used to implement hello.HelloService
type server struct {
	pb.UnimplementedHelloServiceServer
}

// SayHello implements hello.HelloService
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	message := fmt.Sprintf("Hello, %s!", req.GetName())
	return &pb.HelloResponse{Message: message}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8201") // gRPC server on port 8201
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register the HelloService with gRPC
	pb.RegisterHelloServiceServer(grpcServer, &server{})

	// Register the gRPC health service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	// Set the health status to SERVING for the default service (empty string)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	fmt.Println("gRPC server is running on port 8201")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
