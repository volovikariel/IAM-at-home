package models

import (
	"errors"
	"fmt"
)

type User struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("username is required")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

type UserResponse struct {
	Name string `json:"username"`
	Rel  Rel    `json:"rel"`
}

type Rel struct {
	Self string `json:"self"`
}

type UserStore interface {
	// Returns an error if user already exists
	Add(username string, password string) error
	// Returns an error if user doesn't exist
	Get(username string) (*User, error)
	Set(username string, password string) error
	Delete(username string) error
}

type InMemoryUserStore struct {
	UserStore

	users []User
}

func (i *InMemoryUserStore) Add(username string, password string) error {
	i.users = append(i.users, User{Name: username, Password: password})
	return nil
}

func (u *InMemoryUserStore) Get(username string) (*User, error) {
	for _, user := range u.users {
		if user.Name == username {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("User %q not found", username)
}

func (i *InMemoryUserStore) Set(username string, password string) error {
	for _, user := range i.users {
		if user.Name == username {
			user.Password = password
			return nil
		}
	}
	return fmt.Errorf("User %q not found", username)
}

func (i *InMemoryUserStore) Delete(username string) error {
	for index, user := range i.users {
		if user.Name == username {
			i.users = append(i.users[:index], i.users[index+1:]...)
			return nil
		}
	}
	return fmt.Errorf("User %q not found", username)
}

func NewInMemoryUserStore(users []User) *InMemoryUserStore {
	return &InMemoryUserStore{users: users}
}
