package domain

import (
	"architectSocial/domain/helper"
	"strconv"
)

type RegisterUserDto struct {
	FirstName string
	LastName  string
	Age       string
	Gender    string
	Interests string
	City      string
}

type RegisterUserService func(dto *RegisterUserDto) (*helper.ValidationResult, error)

func registerUserValidateDto(dto *RegisterUserDto) *helper.ValidationResult {
	result := helper.NewValidationResult()

	if dto.FirstName == "" {
		result.AddError("firstName", "First name is empty")
	}
	if dto.LastName == "" {
		result.AddError("firstName", "Last name is empty")
	}

	return result
}

func CreateRegisterUserService(userModel UserRepository) RegisterUserService {
	return func(dto *RegisterUserDto) (*helper.ValidationResult, error) {
		validationResult := registerUserValidateDto(dto)
		if !validationResult.IsValid() {
			return validationResult, nil
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
		userModel.CreateUser(dto.FirstName, dto.LastName, uint8(ageNum), gender, dto.Interests, dto.City)

		return nil, nil
	}
}
