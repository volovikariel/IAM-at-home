package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/volovikariel/IdentityManager/internal/server/gateway/sessions"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/users"
	v1 "github.com/volovikariel/IdentityManager/internal/server/gateway/v1"
)

const (
	DEFAULT_HOST             = "localhost"
	DEFAULT_PORT             = "10000"
	DEFAULT_USERNAME_MIN_LEN = 3
	DEFAULT_USERNAME_MAX_LEN = 20
	DEFAULT_PASSWORD_MIN_LEN = 8
	DEFAULT_PASSWORD_MAX_LEN = 256
)

type ServerConfig struct {
	Host           string
	Port           string
	MinUsernameLen int
	MaxUsernameLen int
	MinPasswordLen int
	MaxPasswordLen int
}

func NewServerConfig() ServerConfig {
	return ServerConfig{
		Host:           DEFAULT_HOST,
		Port:           DEFAULT_PORT,
		MinUsernameLen: DEFAULT_USERNAME_MIN_LEN,
		MaxUsernameLen: DEFAULT_USERNAME_MAX_LEN,
		MinPasswordLen: DEFAULT_PASSWORD_MIN_LEN,
		MaxPasswordLen: DEFAULT_PASSWORD_MAX_LEN,
	}
}

type inMemoryUserStore struct {
	users []users.User
}

func (i *inMemoryUserStore) Add(username string, password string) error {
	i.users = append(i.users, users.User{Name: username, Password: password})
	return nil
}

func (u *inMemoryUserStore) Get(username string) error {
	for _, user := range u.users {
		if user.Name == username {
			return nil
		}
	}
	return fmt.Errorf("User %q found", username)
}

type inMemorySessionStore struct {
	sessions []sessions.Session
}

func (i *inMemorySessionStore) Add(username string, token string) error {
	i.sessions = append(i.sessions, sessions.Session{Username: username, Token: token})
	return nil
}

func main() {
	serverConfig := NewServerConfig()

	port := flag.String("p", "", "Port to listen on")
	host := flag.String("h", "", "Host to listen on")
	flag.Parse()
	if port == nil || *port == "" {
		log.Printf("Port not set, defaulting to %q\n", serverConfig.Port)
	} else {
		serverConfig.Port = *port
	}
	if host == nil || *host == "" {
		log.Printf("Host not set, defaulting to %q\n", serverConfig.Host)
	} else {
		serverConfig.Host = *host
	}

	v1Handler := v1.NewHandler()
	http.Handle("/v1/", v1Handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	serverUrl := fmt.Sprintf("%s:%s", serverConfig.Host, serverConfig.Port)
	log.Printf("Listening on %s\n", serverUrl)
	log.Fatal(http.ListenAndServe(serverUrl, nil))
}
