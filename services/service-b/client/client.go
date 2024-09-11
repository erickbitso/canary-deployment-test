package main

import (
	"context"
	"log"
	"time"

	pb "service-b/service-b/hello" // Adjust this path according to your module

	"google.golang.org/grpc"
)

func main() {
	// Establish a connection to the gRPC server
	conn, err := grpc.Dial("localhost:8021", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new HelloService client
	client := pb.NewHelloServiceClient(conn)

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a request message
	req := &pb.HelloRequest{Name: "World"}

	// Send the request and receive a response
	res, err := client.SayHello(ctx, req)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}

	// Print the response message
	log.Printf("Greeting: %s", res.GetMessage())
}
