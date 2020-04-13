package service

import (
	"architectSocial/model"
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

type RegisterUserService func(dto *RegisterUserDto) (*ValidationResult, error)

func registerUserValidateDto(dto *RegisterUserDto) *ValidationResult {
	result := NewValidationResult()

	if dto.FirstName == "" {
		result.AddError("firstName", "First name is empty")
	}
	if dto.LastName == "" {
		result.AddError("firstName", "Last name is empty")
	}

	return result
}

func CreateRegisterUserService(userModel model.UserModel) RegisterUserService {
	return func(dto *RegisterUserDto) (*ValidationResult, error) {
		validationResult := registerUserValidateDto(dto)
		if !validationResult.IsValid() {
			return validationResult, nil
		}
		ageNum, _ := strconv.Atoi(dto.Age)
		var gender model.UserGender
		if dto.Gender == "male" {
			gender = model.Male
		} else if dto.Gender == "female" {
			gender = model.Female
		} else {
			gender = model.Other
		}
		userModel.CreateUser(dto.FirstName, dto.LastName, uint8(ageNum), gender, dto.Interests, dto.City)

		return nil, nil
	}
}
