package helpers

import (
	"database/sql"
	"fmt"
)

type UserStatisticProvider struct {
	db *sql.DB
}

type PopularCity struct {
	Name  string
	Count int
}
type PopularFirstName struct {
	Name  string
	Count int
}

func CreateUserStatisticProvider(db *sql.DB) *UserStatisticProvider {
	return &UserStatisticProvider{db: db}
}

func (provider *UserStatisticProvider) GetMostPopularCities(count int) ([]PopularCity, error) {
	query := "SELECT stat.city,stat.count FROM (SELECT city, count(1) as count FROM users group by city) as stat ORDER BY count DESC LIMIT ?"
	stmt, err := provider.db.Prepare(query)
	if err != nil {
		return []PopularCity{}, fmt.Errorf("failed to fetch popular cities: %s", err.Error())
	}
	rows, err := stmt.Query(count)
	if err != nil {
		return []PopularCity{}, fmt.Errorf("failed to fetch popular cities: %s", err.Error())
	}

	var result []PopularCity
	for rows.Next() {
		item := PopularCity{}
		if err := rows.Scan(&item.Name, &item.Count); err != nil {
			return []PopularCity{}, fmt.Errorf("failed to fetch popular cities: %s", err.Error())
		}

		result = append(result, item)
	}

	return result, nil
}

func (provider *UserStatisticProvider) GetMostPopularFirstNames(count int) ([]PopularFirstName, error) {
	query := "SELECT stat.first_name,stat.count FROM (SELECT first_name, count(1) as count FROM users group by first_name) as stat ORDER BY count DESC LIMIT ?"
	stmt, err := provider.db.Prepare(query)
	if err != nil {
		return []PopularFirstName{}, fmt.Errorf("failed to fetch popular cities: %s", err.Error())
	}
	rows, err := stmt.Query(count)
	if err != nil {
		return []PopularFirstName{}, fmt.Errorf("failed to fetch popular cities: %s", err.Error())
	}

	var result []PopularFirstName
	for rows.Next() {
		item := PopularFirstName{}
		if err := rows.Scan(&item.Name, &item.Count); err != nil {
			return []PopularFirstName{}, fmt.Errorf("failed to fetch popular cities: %s", err.Error())
		}

		result = append(result, item)
	}

	return result, nil
}
