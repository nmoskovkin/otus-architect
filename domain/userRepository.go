package domain

import "github.com/google/uuid"

type UserRepository interface {
	CreateUser(id uuid.UUID, firstName string, lastName string, age uint8, gender UserGender, interests string, city string, password string) error
}
