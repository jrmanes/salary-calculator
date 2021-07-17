package server

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	server *http.Server
}

// New initialize a new server with configuration.
func New(port string) (*Server, error) {
	r := chi.NewRouter()

	// A good base middleware stack
	// middleware which show us logs from the request to our API
	r.Use(middleware.Logger)
	// It allows our server to recover from possible panics
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Use(middleware.SetHeader("AllowCredentials", "true"))
	r.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))
	r.Use(middleware.SetHeader("Access-Control-Allow-Methods", "DELETE, POST, PUT, GET, OPTIONS"))
	r.Use(middleware.SetHeader("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With"))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Mount the admin sub-router
	r.Mount("/api/v1", createRouter())

	// validate the port received, if not, use default 8080
	if len(port) == 0 {
		port = "8080"
	}

	serv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server := Server{server: serv}

	return &server, nil
}

// Start the server.
func (serv *Server) Start() {
	log.Printf("Server running on http://localhost%s", serv.server.Addr)
	log.Fatal(serv.server.ListenAndServe())
}
