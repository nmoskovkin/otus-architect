package domain

import (
	"architectSocial/domain/helper"
	"fmt"
	"github.com/google/uuid"
)

type AuthUserDto struct {
	Login    string
	Password string
}

type AuthUserService func(dto *AuthUserDto) (*helper.ValidationResult, bool, error)

func authUserValidateDto(dto *AuthUserDto) *helper.ValidationResult {
	result := helper.NewValidationResult()
	if dto.Login == "" {
		result.AddError("Login", "Login is empty")
	}
	if dto.Login != "" {
		_, err := uuid.Parse(dto.Login)
		if err != nil {
			result.AddError("Login", "Login must be correct uuid")
		}
	}

	if dto.Password == "" {
		result.AddError("Password", "Password is empty")
	}

	return result
}

func CreateAuthUserService(userRepository UserRepository) AuthUserService {
	return func(dto *AuthUserDto) (*helper.ValidationResult, bool, error) {
		validationResult := authUserValidateDto(dto)
		if !validationResult.IsValid() {
			return validationResult, false, nil
		}

		uuid_, err := uuid.Parse(dto.Login)
		if err != nil {
			return nil, false, fmt.Errorf("failed to parse Login, error: %s", err.Error())
		}
		isAuth, err := userRepository.ExistsWithIdAndPassword(uuid_, dto.Password)
		if err != nil {
			return nil, false, fmt.Errorf("failed find user, error: %s", err.Error())
		}
		return nil, isAuth, nil
	}
}
