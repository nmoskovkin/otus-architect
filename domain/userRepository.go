package domain

type UserRepository interface {
	CreateUser(firstName string, lastName string, age uint8, gender UserGender, interests string, city string)
}
