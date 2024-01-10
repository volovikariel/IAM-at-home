package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/config"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/models"
)

type UserHandler struct {
	http.Handler
	userStore    models.UserStore
	sessionStore models.SessionStore
}

func (u *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remainingPath := strings.TrimPrefix(r.RequestURI, "/v1/users")
	pathParameters := strings.Split(remainingPath, "/")
	if remainingPath == "" {
		// /v1/users
		switch r.Method {
		case http.MethodPost:
			u.CreateUser(w, r)
		default:
			http.Error(w, fmt.Sprintf("Unsupported method: %s for /v1/users", r.Method), http.StatusMethodNotAllowed)
		}
	} else if len(pathParameters) == 1 {
		// /v1/users/{username}
		username := pathParameters[0]
		switch r.Method {
		case http.MethodGet:
			u.GetUser(w, r, username)
		case http.MethodPatch:
			u.UpdateUser(w, r, username)
		case http.MethodDelete:
			u.DeleteUser(w, r, username)
		default:
			http.Error(w, fmt.Sprintf("Unsupported method: %s for /v1/users/{username}", r.Method), http.StatusMethodNotAllowed)
		}

	} else if len(pathParameters) == 2 && pathParameters[1] == "sessions" {
		// /v1/users/sessions/{username}
		NewSessionsHandler(u.sessionStore, u.userStore).ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var user models.User
	err := decoder.Decode(&user)
	if err != nil {
		log.Printf("Could not parse request body: %v\n", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	username := user.Name
	password := user.Password

	if isValid := ValidStringLength(username, config.MinUsernameLength, config.MaxUsernameLength); !isValid {
		// Username is too short or too long
		http.Error(w, "Username is too short or too long", http.StatusBadRequest)
		return
	}
	if _, err := u.userStore.Get(username); err == nil {
		// User already exists
		w.WriteHeader(http.StatusConflict)
		return
	}
	if isValid := ValidStringLength(password, config.MinPasswordLength, config.MaxPasswordLength); !isValid {
		// Password is too short or too long
		http.Error(w, "Password is too short or too long", http.StatusBadRequest)
		return
	}
	if err := u.userStore.Add(username, password); err != nil {
		// Create user failed
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userResponse := models.UserResponse{
		Name: username,
		Rel:  models.Rel{Self: "/v1/users/" + username},
	}
	ur, err := json.Marshal(userResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(ur)
}

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, username string) {
	if isValid := ValidStringLength(username, config.MinUsernameLength, config.MaxUsernameLength); !isValid {
		http.Error(w, "Username is too short or too long", http.StatusBadRequest)
		return
	}
	_, err := u.userStore.Get(username)
	if err != nil {
		// user doesn't exist
		w.WriteHeader(http.StatusNotFound)
		return
	}
	userResponse := models.UserResponse{
		Name: username,
		Rel:  models.Rel{Self: "/v1/users/" + username},
	}
	ur, err := json.Marshal(userResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(ur)
}

func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, username string) {
	if isValid := ValidStringLength(username, config.MinUsernameLength, config.MaxUsernameLength); !isValid {
		http.Error(w, "Username is too short or too long", http.StatusBadRequest)
		return
	}
	_, err := u.userStore.Get(username)
	if err != nil {
		// user doesn't exist
		// TODO: ambiguous whether user is not found or session is not found
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	decoder := json.NewDecoder(r.Body)
	// ensure body body contains ONLY the password & session fields
	decoder.DisallowUnknownFields()
	var updateUserRequest models.UpdateUserRequest
	err = decoder.Decode(&updateUserRequest)
	if err != nil {
		log.Printf("Could not parse request body: %v\n", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// extract password & session (as we'll be updating the user's password)
	password := updateUserRequest.Password
	sessionToken := updateUserRequest.Session
	// TODO: Verify that this account is editable by this session token owner
	accountEditable := true
	if !accountEditable {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	session, _ := u.sessionStore.Get(username)
	// session != nil to verify that the account has a sessiontoken in the first place
	sessionTokenOk := session != nil && sessionToken == session.Token
	if !sessionTokenOk {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("WWW_Authenticate: Bearer charset=\"UTF-8\""))
		return
	}
	// update user's password
	if err := u.userStore.Set(username, password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userResponse := models.UserResponse{
		Name: username,
		Rel:  models.Rel{Self: "/v1/users/" + username},
	}
	ur, err := json.Marshal(userResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(ur)
}

func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, username string) {
	if isValid := ValidStringLength(username, config.MinUsernameLength, config.MaxUsernameLength); !isValid {
		http.Error(w, "Username is too short or too long", http.StatusBadRequest)
		return
	}
	// check if the header params contains ONLY the session-token
	var sessionToken string
	for k, v := range r.Header {
		if k != "session-token" || len(v) != 1 {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		} else {
			sessionToken = v[0]
		}
	}
	// check if user exists
	_, err := u.userStore.Get(username)
	if err != nil {
		// user does not exist
		// no changes were made
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// TODO: Verify that this account is editable by this session token owner
	accountEditable := true
	if !accountEditable {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	session, _ := u.sessionStore.Get(username)
	// session != nil to verify that the account has a sessiontoken in the first place
	sessionTokenOk := session != nil && sessionToken == session.Token
	if !sessionTokenOk {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("WWW_Authenticate: Bearer charset=\"UTF-8\""))
		return
	}

	// delete user
	if err := u.userStore.Delete(username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func NewUsersHandler(userStore models.UserStore, sessionStore models.SessionStore) *UserHandler {
	return &UserHandler{userStore: userStore, sessionStore: sessionStore}
}

func ValidStringLength(str string, minLength int, maxLength int) bool {
	if len(str) < minLength || len(str) > maxLength {
		return false
	}
	return true
}
