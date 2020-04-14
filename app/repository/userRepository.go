package repository

import (
	"architectSocial/domain"
	"database/sql"
)

type MysqlUserModel struct {
	db *sql.DB
}

func CreateUserRepository(db *sql.DB) *MysqlUserModel {
	return &MysqlUserModel{db: db}
}

func (model *MysqlUserModel) CreateUser(firstName string, lastName string, age uint8, gender domain.UserGender, interests string, city string) {

}
