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

type UserItem struct {
	Id        string
	FirstName string
	LastName  string
	Age       uint8
	Interests string
	City      string
	Gender    uint8
}

type ListFilter struct {
	Query string
}

type ListParams struct {
	Filter ListFilter
	Offset int
	Limit  int
}

func (params *ListFilter) Validate() error {
	if len([]rune(params.Query)) > 15 {
		return errors.New("invalid query")
	}
	return nil
}

func (params *ListParams) Validate() error {
	err := params.Filter.Validate()
	if err != nil {
		return fmt.Errorf("invalid filter: %s", err.Error())
	}
	if params.Limit < 0 {
		return fmt.Errorf("invalid limit")
	}
	if params.Offset < 0 {
		return fmt.Errorf("invalid from")
	}

	return nil
}

func CreateListParams() *ListParams {
	return &ListParams{
		Filter: ListFilter{},
		Offset: 0,
		Limit:  15,
	}
}

func CreateMysqlUserRepository(db *sql.DB) *MysqlUserRepository {
	return &MysqlUserRepository{db: db}
}

func (repository *MysqlUserRepository) Create(id uuid.UUID, login string, firstName string, lastName string, age uint8, gender domain.UserGender, interests string, city string, password string) error {
	stmt, err := repository.db.Prepare("INSERT INTO users (id, login, first_name, last_name, age,  gender, interests, city, salt, password) VALUES (?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return errors.New("failed to create user, error: " + err.Error())
	}

	salt := helpers.RandString(16)
	hash, err := bcrypt.GenerateFromPassword([]byte(password+salt), 10)
	if err != nil {
		return errors.New("failed to create user, error: " + err.Error())
	}

	_, err = stmt.Exec(id.String(), login, firstName, lastName, age, gender, interests, city, salt, hash)
	if err != nil {
		return errors.New("failed to create user, error: " + err.Error())
	}

	return nil
}

func (repository *MysqlUserRepository) CreateMany(items []domain.CreateManyItem) error {
	if len(items) == 0 {
		return nil
	}

	sql_ := "INSERT INTO users (id, login, first_name, last_name, age,  gender, interests, city, salt, password) VALUES "
	args := []interface{}{}
	for i, item := range items {
		if i > 0 {
			sql_ += ","
		}
		sql_ += "(?,?,?,?,?,?,?,?,?,?)"
		salt := helpers.RandString(16)
		hash, err := bcrypt.GenerateFromPassword([]byte(item.Password+salt), 1)
		if err != nil {
			return errors.New("failed to create user, error: " + err.Error())
		}
		args = append(args, item.Id.String(), item.Login, item.FirstName, item.LastName, item.Age, item.Gender, item.Interests, item.City, salt, hash)
	}

	stmt, err := repository.db.Prepare(sql_)
	if err != nil {
		panic(err)
		return errors.New("failed to create user, error: " + err.Error())
	}

	_, err = stmt.Exec(args...)
	if err != nil {
		panic(err)
		return errors.New("failed to create user, error: " + err.Error())
	}

	return nil
}

func (repository *MysqlUserRepository) ExistsWithLoginAndPassword(login string, password string) (string, error) {
	stmt, err := repository.db.Prepare("SELECT id,password,salt FROM users WHERE login=?")
	if err != nil {
		return "", fmt.Errorf("failed to fetch user, error: %s", err.Error())
	}
	rows, err := stmt.Query(login)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user, error: %s", err.Error())
	}
	if !rows.Next() {
		return "", domain.UserErrNotFound
	}
	var dbPassword []byte
	var dbSalt []byte
	var id []byte
	if err := rows.Scan(&id, &dbPassword, &dbSalt); err != nil {
		return "", fmt.Errorf("failed to fetch user, error: %s", err.Error())
	}
	err = bcrypt.CompareHashAndPassword(dbPassword, append([]byte(password), dbSalt...))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return "", domain.UserErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("failed to compare password, error: %s", err.Error())
	}

	return string(id), nil
}

func (repository *MysqlUserRepository) GetByIds(ids []string) ([]UserItem, error) {
	wherePart := "WHERE id IN ("
	var args []interface{}
	for i, id := range ids {
		if i > 0 {
			wherePart += ","
		}
		wherePart += "?"
		args = append(args, id)
	}
	wherePart += ")"

	stmt, err := repository.db.Prepare("SELECT id,first_name,last_name,age,interests,city,gender FROM users " + wherePart + " ORDER BY id")
	if err != nil {
		return []UserItem{}, errors.New("failed to fetch user, error: " + err.Error())
	}
	args = append(args)
	rows, err := stmt.Query(args...)
	if err != nil {
		return []UserItem{}, errors.New("failed to fetch user, error: " + err.Error())
	}

	var result []UserItem
	for rows.Next() {
		item := UserItem{}
		if err := rows.Scan(&item.Id, &item.FirstName, &item.LastName, &item.Age, &item.Interests, &item.City, &item.Gender); err != nil {
			return []UserItem{}, err
		}

		result = append(result, item)
	}

	return result, nil
}

func (repository *MysqlUserRepository) GetAll(params *ListParams) ([]UserItem, error) {
	err := params.Validate()
	if err != nil {
		return []UserItem{}, fmt.Errorf("invalid params: %s", err.Error())
	}
	var query string
	var args []interface{}
	if len([]rune(params.Filter.Query)) == 0 {
		query = "SELECT id,first_name,last_name,age,interests,city,gender FROM users ORDER BY id LIMIT ?,?"
		args = append(args, params.Offset, params.Limit)
	} else {
		query = "SELECT id,first_name,last_name,age,interests,city,gender FROM users WHERE first_name like ? AND last_name like ? ORDER BY id LIMIT ?"
		args = append(args, params.Filter.Query+"%", params.Filter.Query+"%", params.Limit)
	}
	stmt, err := repository.db.Prepare(query)
	if err != nil {
		return []UserItem{}, errors.New("failed to fetch user, error: " + err.Error())
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		return []UserItem{}, errors.New("failed to fetch user, error: " + err.Error())
	}

	var result []UserItem
	for rows.Next() {
		item := UserItem{}
		if err := rows.Scan(&item.Id, &item.FirstName, &item.LastName, &item.Age, &item.Interests, &item.City, &item.Gender); err != nil {
			return []UserItem{}, err
		}

		result = append(result, item)
	}

	return result, nil
}
