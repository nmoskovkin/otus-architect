package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port       uint16
	MysqlDSN   string
	SessionKey string
}

func NewConfig() (*Config, error) {
	mysqlDsn, err := getEnv("SOCIAL_MYSQL_DSN")
	if err != nil {
		return nil, errors.New("failed to load config, error: " + err.Error())
	}
	port, err := getEnvAsInt("SOCIAL_PORT")
	if err != nil {
		return nil, errors.New("failed to load config, error: " + err.Error())
	}
	if port < 0 {
		return nil, errors.New("failed to load config, port is less than zero")
	}
	sessionKey, err := getEnv("SOCIAL_SESSION_KEY")
	if err != nil {
		return nil, errors.New("failed to load config, error: " + err.Error())
	}
	return &Config{
		Port:       uint16(port),
		MysqlDSN:   mysqlDsn,
		SessionKey: sessionKey,
	}, nil
}

func getEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}

	return "", errors.New("undefined variable: \"" + key + "\"")
}

func getEnvAsInt(name string) (int, error) {
	valueStr, err := getEnv(name)
	if err != nil {
		return 0, err
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, nil
	}

	return value, nil
}

func getEnvAsBool(name string) (bool, error) {
	valStr, err := getEnv(name)
	if err != nil {
		return false, err
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return false, err
	}

	return val, nil
}

func getEnvAsSlice(name string, sep string) ([]string, error) {
	valStr, err := getEnv(name)
	if err != nil {
		return []string{}, err
	}

	if valStr == "" {
		return []string{}, nil
	}

	return strings.Split(valStr, sep), nil
}
