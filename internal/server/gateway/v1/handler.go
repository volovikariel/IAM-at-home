package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/volovikariel/IdentityManager/internal/server/gateway/sessions"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/users"
)

type v1Handler struct {
	userStore    users.UserStore
	sessionStore sessions.SessionStore
}

func NewHandler() http.Handler {
	return v1Handler{}
}

func (h v1Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remainingPath := strings.TrimPrefix(r.RequestURI, "/v1")
	pathParameters := strings.Split(remainingPath, "/")

	// /v1/ not defined
	if len(pathParameters) == 0 {
		http.NotFound(w, r)
		return
	}
	// /v1/users
	if pathParameters[0] != "users" {
		http.NotFound(w, r)
		return
	}

	userHandler := users.NewHandler(h.userStore)
	sessionHandler := sessions.NewHandler(h.sessionStore, h.userStore)
	// /v1/users/
	if len(pathParameters) == 1 {
		switch r.Method {
		case http.MethodPost:
			userHandler.CreateUser(w, r)
		default:
			http.Error(w, fmt.Sprintf("Unsupported method: %s for /v1/users", r.Method), http.StatusMethodNotAllowed)
		}
	}

	// /v1/users/{username}
	if len(pathParameters) == 2 {
		// TODO: extract username from path
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUser(w, r)
		case http.MethodPatch:
			userHandler.UpdateUser(w, r)
		case http.MethodDelete:
			userHandler.DeleteUser(w, r)
		default:
			http.Error(w, fmt.Sprintf("Unsupported method: %s for /v1/users/{username}", r.Method), http.StatusMethodNotAllowed)
		}
	}

	// /v1/users/sessions/{username}
	if len(pathParameters) == 3 {
		// TODO: extract username from path
		switch r.Method {
		case http.MethodPost:
			sessionHandler.CreateSession(w, r)
		case http.MethodDelete:
			sessionHandler.TerminateSession(w, r)
		default:
			http.Error(w, fmt.Sprintf("Unsupported method: %s for /v1/users/sessions/{username}", r.Method), http.StatusMethodNotAllowed)
		}
	}

	http.NotFound(w, r)
}
