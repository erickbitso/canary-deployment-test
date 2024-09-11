package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	//"google.golang.org/grpc"
	//pb "service-a/service-b/hello"
)

const (
	url         = "http://service-b.default.svc.cluster.local:8080/hello"
	tcpAddress  = "service-b.default.svc.cluster.local:8000"
	grpcAddress = "service-b.default.svc.cluster.local:8021"
)

func main() {

	// Create a ticker that triggers every 100 milliseconds
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		sendRequest(url, tcpAddress, grpcAddress)
	}
}

func sendRequest(url, tcpAddress, grpcAddress string) {
	// Send a GET request to the specified URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to send HTTP request: %v\n", err)
	} else {
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed to read HTTP response body: %v\n", err)
		} else {
			// Print the HTTP response status and body
			fmt.Printf("HTTP Response Status: %s, Body: %s\n", resp.Status, string(body))
		}
	}

	// // Send a TCP request to the specified address
	// conn, err := net.Dial("tcp", tcpAddress)
	// if err != nil {
	// 	fmt.Printf("Failed to send TCP request: %v\n", err)
	// } else {
	// 	defer conn.Close()

	// 	// Read the TCP response
	// 	tcpResponse, err := ioutil.ReadAll(conn)
	// 	if err != nil {
	// 		fmt.Printf("Failed to read TCP response: %v\n", err)
	// 	} else {
	// 		// Print the TCP response
	// 		fmt.Printf("TCP Response Status: %s\n", string(tcpResponse))
	// 	}
	// }

	// // Establish a connection to the gRPC server
	// connx, err := grpc.Dial("localhost:8021", grpc.WithInsecure(), grpc.WithBlock())
	// if err != nil {
	// 	log.Fatalf("Did not connect: %v", err)
	// }
	// defer conn.Close()

	// // Create a new HelloService client
	// client := pb.NewHelloServiceClient(connx)

	// // Create a context with a timeout
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()

	// // Create a request message
	// req := &pb.HelloRequest{Name: "World"}

	// // Send the request and receive a response
	// res, err := client.SayHello(ctx, req)
	// if err != nil {
	// 	fmt.Printf("Could not greet: %v\n", err)
	// 	return
	// }

	// // Print the response message
	// fmt.Printf("gRPC Response Status: %s\n", res.GetMessage())
}
