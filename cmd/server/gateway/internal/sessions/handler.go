package sessions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/volovikariel/IdentityManager/cmd/server/gateway/internal/users"
)

type SessionHandler struct {
	sessionStore SessionStore
	userStore    users.UserStore // So we can check if users exist
}

func (s *SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: Sessions")
	switch r.Method {
	case http.MethodPost:
		s.CreateSession(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func NewSessionHandler(sessionStore SessionStore, userStore users.UserStore) *SessionHandler {
	return &SessionHandler{sessionStore: sessionStore, userStore: userStore}
}

func (s *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var session Session

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&session); err != nil {
		w.Write([]byte("Could not parse request body"))
		return
	}
	if err := s.userStore.Get(session.Username); err != nil {
		http.Error(w, fmt.Sprintf("User %q not found, cannot create session", session.Username), http.StatusNotFound)
		return
	}

	log.Printf("Creating session for user %q\n", session.Username)
	s.sessionStore.Add(session.Username, session.Token)
	// Token should be some unique value
	token := "1"
	w.Write([]byte(fmt.Sprintf("Session created for user %q with token %q", session.Username, token)))
}
