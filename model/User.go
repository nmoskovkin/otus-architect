package model

import "database/sql"

type UserGender int

const (
	Male UserGender = iota
	Female
	Other
)

type UserModel interface {
	CreateUser(firstName string, lastName string, age uint8, gender UserGender, interests string, city string)
}

type MysqlUserModel struct {
	db *sql.DB
}

func (model *MysqlUserModel) CreateUser(firstName string, lastName string, age uint8, gender UserGender, interests string, city string) {

}

func CreateMysqlUserModel(db *sql.DB) *MysqlUserModel {
	return &MysqlUserModel{db: db}
}
