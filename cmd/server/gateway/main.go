package main

import (
	"flag"
	"net/http"

	v1 "github.com/volovikariel/IdentityManager/internal/server/gateway/v1"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/config"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/handlers"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/middleware"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/models"
)

func main() {
	serverConfig := &config.Server{}
	flag.StringVar(&serverConfig.Host, "h", config.DEFAULT_HOST, "Host to listen on")
	flag.StringVar(&serverConfig.Port, "p", config.DEFAULT_PORT, "Port to listen on")
	flag.Parse()
	server := v1.NewServer(serverConfig)

	memoryStore := &models.InMemoryUserStore{}
	sessionStore := &models.InMemorySessionStore{}
	v1Handler := handlers.NewV1Handler(memoryStore, sessionStore)
	mux := http.NewServeMux()
	mux.Handle("/v1/", middleware.LoggingMiddleware(v1Handler))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	server.Start(mux)
}
