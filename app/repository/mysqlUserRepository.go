package repository

import (
	"architectSocial/app/helpers"
	"architectSocial/domain"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type MysqlUserRepository struct {
	db *sql.DB
}

type FindAllItem struct {
	Id        string
	FirstName string
	LastName  string
	Age       uint8
	Interests string
	City      string
	Gender    uint8
}

func CreateMysqlUserRepository(db *sql.DB) *MysqlUserRepository {
	return &MysqlUserRepository{db: db}
}

func (model *MysqlUserRepository) Create(id uuid.UUID, firstName string, lastName string, age uint8, gender domain.UserGender, interests string, city string, password string) error {
	stmt, err := model.db.Prepare("INSERT INTO users (id, first_name, last_name, age,  gender, interests, city, salt, password) VALUES (?,?,?,?,?,?,?,?,?)")

	if err != nil {
		return errors.New("failed to create user, error: " + err.Error())
	}

	salt := helpers.RandString(16)
	hash, err := bcrypt.GenerateFromPassword([]byte(password+salt), 10)
	if err != nil {
		return errors.New("failed to create user, error: " + err.Error())
	}

	_, err = stmt.Exec(id.String(), firstName, lastName, age, gender, interests, city, salt, hash)
	if err != nil {
		return errors.New("failed to create user, error: " + err.Error())
	}

	return nil
}

func (model *MysqlUserRepository) ExistsWithIdAndPassword(id uuid.UUID, password string) (bool, error) {
	stmt, err := model.db.Prepare("SELECT password,salt FROM users WHERE id=?")
	if err != nil {
		return false, fmt.Errorf("failed to fetch user, error: %s", err.Error())
	}
	rows, err := stmt.Query(id.String())
	if err != nil {
		return false, fmt.Errorf("failed to fetch user, error: %s", err.Error())
	}
	if !rows.Next() {
		return false, nil
	}
	var dbPassword []byte
	var dbSalt []byte
	if err := rows.Scan(&dbPassword, &dbSalt); err != nil {
		return false, fmt.Errorf("failed to fetch user, error: %s", err.Error())
	}
	err = bcrypt.CompareHashAndPassword(dbPassword, append([]byte(password), dbSalt...))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to compare password, error: %s", err.Error())
	}

	return true, nil
}

func (model *MysqlUserRepository) GetAll() ([]FindAllItem, error) {
	stmt, err := model.db.Prepare("SELECT id,first_name,last_name,age,interests,city,gender FROM users")
	if err != nil {
		return []FindAllItem{}, errors.New("failed to fetch user, error: " + err.Error())
	}
	rows, err := stmt.Query()
	if err != nil {
		return []FindAllItem{}, errors.New("failed to fetch user, error: " + err.Error())
	}

	var result []FindAllItem
	for rows.Next() {
		item := FindAllItem{}
		if err := rows.Scan(&item.Id, &item.FirstName, &item.LastName, &item.Age, &item.Interests, &item.City, &item.Gender); err != nil {
			return []FindAllItem{}, err
		}

		result = append(result, item)
	}

	return result, nil
}
