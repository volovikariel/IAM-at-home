package v1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/config"
)

type Server struct {
	config *config.Server
}

func NewServer(config *config.Server) *Server {
	return &Server{config: config}
}

func (s *Server) Start(mux *http.ServeMux) {
	serverUrl := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	log.Printf("Listening on %s\n", serverUrl)
	log.Fatal(http.ListenAndServe(serverUrl, mux))
}
