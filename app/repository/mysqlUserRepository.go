package repository

import (
	"architectSocial/app/helpers"
	"architectSocial/domain"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
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

type GetAllFilter struct {
	Id          string
	Ids         []string
	FilterByIds bool
	Query       string
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
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	fmt.Printf("Start: %s\n", formatted)
	if len(items) == 0 {
		return nil
	}

	//fmt.Println(items)
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

	//fmt.Println(sql)
	//fmt.Println(sql_)
	stmt, err := repository.db.Prepare(sql_)
	//fmt.Println(sql_)

	if err != nil {
		panic(err)
		return errors.New("failed to create user, error: " + err.Error())
	}

	_, err = stmt.Exec(args...)
	//fmt.Println(err)

	if err != nil {
		panic(err)
		return errors.New("failed to create user, error: " + err.Error())
	}
	t = time.Now()
	formatted = fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	fmt.Println(formatted)
	fmt.Printf("End: %s\n", formatted)

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

func (repository *MysqlUserRepository) GetAll(filter GetAllFilter, from int, count int) ([]FindAllItem, error) {
	wherePart := ""
	args := []interface{}{}
	if filter.Id != "" {
		// Query builder :(
		args = append(args, filter.Id)
		wherePart = "WHERE id=?"
	} else if len(filter.Ids) > 0 {
		wherePart = "WHERE id IN ("
		for i, id := range filter.Ids {
			if i > 0 {
				wherePart += ","
			}
			wherePart += "?"
			args = append(args, id)
		}
		wherePart += ")"
	} else if filter.FilterByIds && len(filter.Ids) == 0 {
		return []FindAllItem{}, nil
	} else if filter.Query != "" {
		wherePart = "WHERE first_name LIKE ? OR last_name LIKE ?"
		args = append(args, filter.Query+"%", filter.Query+"%")
	}

	stmt, err := repository.db.Prepare("SELECT id,first_name,last_name,age,interests,city,gender FROM users " + wherePart + " ORDER BY id LIMIT ?,?")
	if err != nil {
		return []FindAllItem{}, errors.New("failed to fetch user, error: " + err.Error())
	}

	args = append(args, from, count)
	rows, err := stmt.Query(args...)
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
