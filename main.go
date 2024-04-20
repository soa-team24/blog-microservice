package main

import (
	"blog-microservice/handler"
	"blog-microservice/repo"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

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

	//Initialize the router and add a middleware for all the requests
	router := mux.NewRouter()
	//kada istanciramo router, preko .Use treba da mu prosledimo Middleware(fju za kreiranje middleware-a)
	router.Use(blogHandelr.MiddlewareContentTypeSet) //osnovni middleware

	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/blog", blogHandelr.GetAllBlogs)

	getByIdRouter := router.Methods(http.MethodGet).Subrouter()
	getByIdRouter.HandleFunc("/blog/{id}", blogHandelr.GetBlogById)

	postRouter := router.Methods(http.MethodPost).Subrouter() //pravim novu instancu routera na osnovu inicijalnog routera i kazem da ce metod biti post
	postRouter.HandleFunc("/blog", blogHandelr.PostBlog)      // kazem koji hendler ce hendlati taj zahtev
	postRouter.Use(blogHandelr.MiddlewareBlogDeserialization) //njega presrece ovaj middleware

	getByAuthorNameRouter := router.Methods(http.MethodGet).Subrouter()
	getByAuthorNameRouter.HandleFunc("/blog/byUser/{userId}", blogHandelr.GetBlogsByAuthorId)

	updateRouter := router.Methods(http.MethodPatch).Subrouter()
	updateRouter.HandleFunc("/blog/{id}", blogHandelr.UpdateBlog)
	updateRouter.Use(blogHandelr.MiddlewareBlogDeserialization)

	deleteRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/blog/{id}", blogHandelr.DeleteBlog)

	addVoteRouter := router.Methods(http.MethodPatch).Subrouter()
	addVoteRouter.HandleFunc("/blog/votes/{id}", blogHandelr.AddVote)

	changeVoteRouter := router.Methods(http.MethodPatch).Subrouter()
	changeVoteRouter.HandleFunc("/blog/votes/{id}", blogHandelr.ChangeVote)

	allowedOrigins := handlers.AllowedOrigins([]string{"*"}) // Allow all origins
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{
		"Content-Type",
		"Authorization",
		"X-Custom-Header",
	})

	cors := handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)

	//Initialize the server
	server := http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Println("Server listening on port", port)
	//Distribute all the connections to goroutines
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	//Try to shutdown gracefully
	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")
}
