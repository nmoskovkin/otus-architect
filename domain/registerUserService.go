package domain

import (
	"architectSocial/domain/helper"
	"errors"
	"github.com/google/uuid"
	"strconv"
)

type RegisterUserDto struct {
	FirstName            string
	LastName             string
	Age                  string
	Gender               string
	Interests            string
	City                 string
	Password             string
	PasswordConfirmation string
}

type RegisterUserService func(dto *RegisterUserDto) (*helper.ValidationResult, string, error)

func registerUserValidateDto(dto *RegisterUserDto) *helper.ValidationResult {
	result := helper.NewValidationResult()

	// TODO maybe validate empty values in a loop
	if dto.FirstName == "" {
		result.AddError("firstName", "First name is empty")
	}
	if dto.LastName == "" {
		result.AddError("firstName", "Last name is empty")
	}
	if dto.Age == "" {
		result.AddError("age", "Age is empty")
	}
	if dto.Gender == "" {
		result.AddError("gender", "Gender is empty")
	}
	if dto.Interests == "" {
		result.AddError("interests", "Interests are empty")
	}
	if dto.City == "" {
		result.AddError("city", "City is empty")
	}
	if dto.Password == "" {
		result.AddError("Password", "Password is empty")
	}
	if dto.PasswordConfirmation == "" {
		result.AddError("passwordConfirm", "Please confirm your Password")
	}
	if dto.Password != "" && dto.PasswordConfirmation != "" && dto.Password != dto.PasswordConfirmation {
		result.AddError("Password", "The passwords don't match")
	}
	if dto.Age != "" {
		num, err := strconv.Atoi(dto.Age)
		if err != nil {
			result.AddError("age", "Age must be a number")
		}
		if err == nil && num <= 0 {
			result.AddError("age", "Age must be greater than zero")
		}
	}

	return result
}

func CreateRegisterUserService(userModel UserRepository) RegisterUserService {
	return func(dto *RegisterUserDto) (*helper.ValidationResult, string, error) {
		validationResult := registerUserValidateDto(dto)
		if !validationResult.IsValid() {
			return validationResult, "", nil
		}
		ageNum, _ := strconv.Atoi(dto.Age)
		var gender UserGender
		if dto.Gender == "male" {
			gender = Male
		} else if dto.Gender == "female" {
			gender = Female
		} else {
			gender = Other
		}
		id, err := uuid.NewUUID()
		if err != nil {
			return nil, "", errors.New("failed to create user, error:" + err.Error())
		}

		err = userModel.Create(id, dto.FirstName, dto.LastName, uint8(ageNum), gender, dto.Interests, dto.City, dto.Password)
		if err != nil {
			return nil, "", errors.New("failed to create user, error:" + err.Error())
		}

		return nil, id.String(), nil
	}
}
