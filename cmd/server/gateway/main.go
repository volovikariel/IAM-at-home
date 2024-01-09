package main

import (
	"flag"
	"fmt"
	"net/http"

	v1 "github.com/volovikariel/IdentityManager/internal/server/gateway/v1"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/config"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/handlers"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/models"
)

type inMemoryUserStore struct {
	models.UserStore

	users []models.User
}

func (i *inMemoryUserStore) Add(username string, password string) error {
	i.users = append(i.users, models.User{Name: username, Password: password})
	return nil
}

func (u *inMemoryUserStore) Get(username string) (*models.User, error) {
	for _, user := range u.users {
		if user.Name == username {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("User %q not found", username)
}

type inMemorySessionStore struct {
	models.SessionStore

	sessions []models.Session
}

func (i *inMemorySessionStore) Add(username string, token string) error {
	i.sessions = append(i.sessions, models.Session{Username: username, Token: token})
	return nil
}

func main() {
	serverConfig := &config.Server{}
	flag.StringVar(&serverConfig.Port, "p", config.DEFAULT_PORT, "Port to listen on")
	flag.StringVar(&serverConfig.Host, "h", config.DEFAULT_HOST, "Host to listen on")
	flag.Parse()
	server := v1.NewServer(serverConfig)

	v1Handler := handlers.NewV1Handler(&inMemoryUserStore{}, &inMemorySessionStore{})
	mux := http.NewServeMux()
	mux.Handle("/v1/", v1Handler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	server.Start(mux)
}
