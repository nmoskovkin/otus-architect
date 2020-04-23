package domain

import (
	"errors"
	"github.com/google/uuid"
)

var UserErrNotFound = errors.New("user not found")

type CreateManyItem struct {
	Id        uuid.UUID
	Login     string
	FirstName string
	LastName  string
	Age       uint8
	Gender    UserGender
	Interests string
	City      string
	Password  string
}

type UserRepository interface {
	Create(id uuid.UUID, login string, firstName string, lastName string, age uint8, gender UserGender, interests string, city string, password string) error
	CreateMany(items []CreateManyItem) error
	ExistsWithLoginAndPassword(login string, password string) (string, error)
}
