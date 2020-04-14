package repository

import (
	"architectSocial/app/helpers"
	"architectSocial/domain"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type MysqlUserModel struct {
	db *sql.DB
}

func CreateUserRepository(db *sql.DB) *MysqlUserModel {
	return &MysqlUserModel{db: db}
}

func (model *MysqlUserModel) CreateUser(id uuid.UUID, firstName string, lastName string, age uint8, gender domain.UserGender, interests string, city string, password string) error {
	stmt, err := model.db.Prepare("INSERT INTO users (id, first_name, last_name, age,  gender, interests, city, salt, password) VALUES (?,?,?,?,?,?,?,?,?)")

	if err != nil {
		return errors.New("failed to create user, error: " + err.Error())
	}

	salt := helpers.RandString(16)
	hash, err := bcrypt.GenerateFromPassword([]byte(salt+password), 10)
	if err != nil {
		return errors.New("failed to create user, error: " + err.Error())
	}

	_, err = stmt.Exec(id.String(), firstName, lastName, age, gender, interests, city, salt, hash)
	if err != nil {
		return errors.New("failed to create user, error: " + err.Error())
	}

	return nil
}
