package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	pb "service-b/service-b/hello" // Adjust this path
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	httpPort = 8080
	tcpPort  = 8000
	grpcPort = 8021
)

type server struct {
	pb.UnimplementedHelloServiceServer
}

func main() {
	// Set up the /hello endpoint
	http.HandleFunc("/hello", helloHandler)

	// Start the HTTP server on port 8080
	go func() {
		fmt.Printf("HTTP server is running on port %d\n", httpPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil); err != nil {
			fmt.Printf("Failed to start HTTP server: %v\n", err)
		}
	}()

	// Start the TCP server on port 8000
	go func() {
		fmt.Printf("TCP server is running on port %d\n", tcpPort)
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", tcpPort))
		if err != nil {
			fmt.Printf("Failed to start TCP server: %v\n", err)
			return
		}
		defer listener.Close()

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("Failed to accept connection: %v\n", err)
				continue
			}
			go handleTCPConnection(conn)
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
		if err != nil {
			fmt.Printf("Failed to listen: %v", err)
			return
		}

		s := grpc.NewServer()
		pb.RegisterHelloServiceServer(s, &server{})
		reflection.Register(s)

		fmt.Printf("gRPC server is running on port %d\n", grpcPort)
		if err := s.Serve(lis); err != nil {
			fmt.Printf("Failed to serve: %v", err)
			return
		}
	}()

	// Prevent the main function from exiting
	select {}
}

// helloHandler responds with "OK"
func helloHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now() // Get the current date and time
	dateTimeStr := now.Format("2006-01-02 15:04:05")
	fmt.Println("Received a request", dateTimeStr)
	fmt.Fprintln(w, "OK", dateTimeStr)
}

// handleTCPConnection handles incoming TCP connections
func handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	message := "OK\n"
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Failed to write to TCP connection: %v\n", err)
	}
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	message := fmt.Sprintf("Hello, %s!", req.GetName())
	return &pb.HelloResponse{Message: message}, nil
}
