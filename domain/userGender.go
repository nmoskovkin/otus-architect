package domain

type UserGender int

const (
	Male UserGender = iota + 1
	Female
	Other
)
