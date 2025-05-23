package server

import (
	"GOrion/internal/logging"
	"GOrion/internal/router"
	"net/http"
	"log"
)

type Server struct {
	Port string
	// For future todo add other server configuration 
	// fields here (e.g., TLSConfig, ReadTimeout)
}

func NewServer(port string) *Server {
	return &Server{
		Port: port,
	}
}

func (s *Server) ServerRun() {
	var serverPort string = ":" + s.Port
	logging.LogAndPrint("Server starting on port %s...", s.Port)

	r := router.SetupRouter()

	log.Fatal(http.ListenAndServe(serverPort, r))
}

func ServerRunOnlyRoutes() {
	router.SetupRouter()
}
