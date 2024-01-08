package sessions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/models"
)

type SessionHandler struct {
	sessionStore models.SessionStore
	userStore    models.UserStore // So we can check if the user exists before creating/terminating their session
}

func NewHandler(sessionStore models.SessionStore, userStore models.UserStore) *SessionHandler {
	return &SessionHandler{sessionStore: sessionStore, userStore: userStore}
}

// TODO: Update this function
func (s *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var session models.Session

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&session); err != nil {
		w.Write([]byte("Could not parse request body"))
		return
	}
	if _, err := s.userStore.Get(session.Username); err != nil {
		http.Error(w, fmt.Sprintf("User %q not found, cannot create session", session.Username), http.StatusNotFound)
		return
	}

	log.Printf("Creating session for user %q\n", session.Username)
	s.sessionStore.Add(session.Username, session.Token)
	// Token should be some unique value
	token := "1"
	w.Write([]byte(fmt.Sprintf("Session created for user %q with token %q", session.Username, token)))
}

// TODO: Update this function
func (s *SessionHandler) TerminateSession(w http.ResponseWriter, r *http.Request) {

}

func (s *SessionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.CreateSession(w, r)
	case http.MethodDelete:
		s.TerminateSession(w, r)
	default:
		http.Error(w, fmt.Sprintf("Unsupported method: %s for /v1/sessions", r.Method), http.StatusMethodNotAllowed)
	}
}
