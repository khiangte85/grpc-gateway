package main

import (
	"context"
	hellopb "grpc-gateway/proto/hello"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type HelloServer struct {
	hellopb.UnimplementedGreeterServer
}

func NewServer() *HelloServer {
	return &HelloServer{}
}

func (*HelloServer) SayHello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloReply, error) {
	name := req.GetName()
	return &hellopb.HelloReply{Message: "Hello " + name}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen", err)
	}

	s := grpc.NewServer()

	hellopb.RegisterGreeterServer(s, &HelloServer{})

	reflection.Register(s)

	go func() {
		log.Println("Serving GRPC at 0.0.0.0:8080")
		log.Fatalln(s.Serve(lis))
	}()

	conn, err := grpc.NewClient(
		"0.0.0.0:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = hellopb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8081",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8081")
	log.Fatalln(gwServer.ListenAndServe())

}
