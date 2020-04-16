package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type MysqlFriendsRepository struct {
	db *sql.DB
}

func CreateMysqlFriendsRepository(db *sql.DB) *MysqlFriendsRepository {
	return &MysqlFriendsRepository{db: db}
}

func (repository *MysqlFriendsRepository) Create(fromUser uuid.UUID, toUser uuid.UUID) error {
	stmt, err := repository.db.Prepare("INSERT INTO friends (from_user, to_user) VALUES (?,?)")

	if err != nil {
		return fmt.Errorf("failed to create friend, error: %s", err.Error())
	}

	_, err = stmt.Exec(fromUser.String(), toUser.String())
	if err != nil {
		return fmt.Errorf("failed to create friend, error: %s", err.Error())
	}
	_, err = stmt.Exec(toUser.String(), fromUser.String())
	if err != nil {
		return fmt.Errorf("failed to create friend, error: %s", err.Error())
	}

	return nil
}

func (repository *MysqlFriendsRepository) GetFriends(fromUser uuid.UUID) ([]string, error) {
	stmt, err := repository.db.Prepare("SELECT to_user FROM friends WHERE from_user=?")

	if err != nil {
		return []string{}, fmt.Errorf("failed to fetch friedns, error: %s", err.Error())
	}

	rows, err := stmt.Query(fromUser.String())
	if err != nil {
		return []string{}, fmt.Errorf("failed to fetch friedns, error: %s", err.Error())
	}
	result := []string{}
	for rows.Next() {
		var toUser string
		err = rows.Scan(&toUser)
		if err != nil {
			return []string{}, fmt.Errorf("failed to fetch data, error: %s", err.Error())
		}
		result = append(result, toUser)
	}
	return result, nil
}

func (repository *MysqlFriendsRepository) AreFriends(user1 uuid.UUID, user2 uuid.UUID) (bool, error) {
	stmt, err := repository.db.Prepare("SELECT * FROM friends WHERE from_user=? and to_user=?")
	if err != nil {
		return false, fmt.Errorf("failed to fetch friedns, error: %s", err.Error())
	}
	rows, err := stmt.Query(user1.String(), user2.String())
	if err != nil {
		return false, fmt.Errorf("failed to fetch friedns, error: %s", err.Error())
	}

	return rows.Next(), nil
}
