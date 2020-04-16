package domain

import (
	"errors"
	"github.com/google/uuid"
)

var UserErrNotFound = errors.New("user not found")

type UserRepository interface {
	Create(id uuid.UUID, login string, firstName string, lastName string, age uint8, gender UserGender, interests string, city string, password string) error
	ExistsWithLoginAndPassword(login string, password string) (string, error)
}
