package api

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/heroku/sheshankkodam-be-todo/pkg/database"
)

type Server struct {
	server 		*http.Server
	db 			*database.Database
	startTime 	time.Time
}

func NewServer(httpPort string, db *database.Database) *Server {
	//corsHandler := handlers.CORS(handlers.IgnoreOptions(), handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}))
	muxRouter := mux.NewRouter()

	methods := handlers.AllowedMethods([]string{"OPTIONS", "DELETE", "GET", "HEAD", "POST", "PUT", "PATCH"})
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Accept", "Access-Control-Allow-Origin"})
	origins := handlers.AllowedOrigins([]string{"*"})
	handler := handlers.CORS(methods, origins, headers)

	server := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        handler(muxRouter),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1000000,
	}
	s := &Server{
		server:server,
		db:db,
		startTime: time.Now().UTC(),
	}

	muxRouter.HandleFunc("/", s.liveHandler).Methods(http.MethodGet)
	muxRouter.HandleFunc("/api/v1/status", s.statusHandler).Methods(http.MethodGet)
	muxRouter.HandleFunc("/api/v1/login", s.loginHandler).Methods(http.MethodPost)
	muxRouter.HandleFunc("/api/v1/customer/{customerId}/add", s.addTodo).Methods(http.MethodPost)
	muxRouter.HandleFunc("/api/v1/customer/{customerId}/todos", s.getTodoTasks).Methods(http.MethodGet)
	muxRouter.HandleFunc("/api/v1/customer/{customerId}/todo/{todoId}", s.deleteTodo).Methods(http.MethodDelete)
	muxRouter.HandleFunc("/api/v1/customer/{customerId}/todo/{todoId}/update", s.updateTodo).Methods(http.MethodPut)

	return s
}

func (s *Server) Run() error {
	l, err := net.Listen("tcp4", s.server.Addr)
	if err != nil {
		return err
	}
	err = s.server.Serve(l)
	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}

func (s *Server) Close() {
	err := s.server.Close()
	if err != nil {
		log.Println("Error closing HTTP server")
	}
}
