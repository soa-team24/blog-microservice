package main

import (
	"blog-microservice/handler"
	"blog-microservice/repo"
	"context"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"soa/grpc/proto/blog"
)

func main() {

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger := log.New(os.Stdout, "[blog-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[blog-store] ", log.LstdFlags)

	// NoSQL: Initialize Blog Repository store
	store, err := repo.New(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.Disconnect(timeoutContext)

	// NoSQL: Checking if the connection was established, proveri da li je sve ok
	store.Ping()

	blogHandelr := handler.NewBlogHandler(logger, store)

	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	blog.RegisterBlogServiceServer(grpcServer, blogHandelr)
	reflection.Register(grpcServer)
	grpcServer.Serve(lis)

	logger.Println("Server stopped")
}
