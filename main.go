package main

import (
	"blog-microservice/handler"
	"blog-microservice/repo"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"blog-microservice/proto/blog"
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

	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	blog.RegisterBlogServiceServer(grpcServer, blogHandelr)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGTERM)

	<-stopCh

	grpcServer.Stop()
}
