package models

import "fmt"

type User struct {
	ID        int    // unique
	Email     string // unique
	Token     string // unique
	Username  string // unique
	Password  string
	Bio       string
	Image     string
	Followers map[int]bool
}

type UserRepository interface {
	NewID() int
	GetByEmail(email string) (res *User, err error)
	GetByUsername(username string) (res *User, err error)
	GetByID(id int) (res *User, err error)
	Add(entity User) error
	UpdateByID(id int, entity User) (err error)
	ListByFollowedBy(userID int) (res []User, err error)
}

type UserByEmailNotFoundError struct {
	Email string
}

func (e UserByEmailNotFoundError) Error() string {
	return fmt.Sprintf("user with email '%s' not found", e.Email)
}

type UserByUsernameNotFoundError struct {
	Username string
}

func (e UserByUsernameNotFoundError) Error() string {
	return fmt.Sprintf("user with username '%s' not found", e.Username)
}

type UserByIDNotFoundError struct {
	ID int
}

func (e UserByIDNotFoundError) Error() string {
	return fmt.Sprintf("user with id '%d' not found", e.ID)
}
