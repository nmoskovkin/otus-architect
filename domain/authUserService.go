package domain

import (
	"architectSocial/domain/helper"
	"fmt"
)

type AuthUserDto struct {
	Login    string
	Password string
}

type AuthUserService func(dto *AuthUserDto, authenticator Authenticator) (*helper.ValidationResult, bool, error)

type Authenticator func(id string)

func authUserValidateDto(dto *AuthUserDto) *helper.ValidationResult {
	result := helper.NewValidationResult()
	if dto.Login == "" {
		result.AddError("Login", "Login is empty")
	}
	if dto.Password == "" {
		result.AddError("Password", "Password is empty")
	}

	return result
}

func CreateAuthUserService(userRepository UserRepository) AuthUserService {
	return func(dto *AuthUserDto, authenticator Authenticator) (*helper.ValidationResult, bool, error) {
		validationResult := authUserValidateDto(dto)
		if !validationResult.IsValid() {
			return validationResult, false, nil
		}
		id, err := userRepository.ExistsWithLoginAndPassword(dto.Login, dto.Password)
		var isAuth bool
		if err == UserErrNotFound {
			isAuth = false
		} else if err != nil {
			return nil, false, fmt.Errorf("failed find user, error: %s", err.Error())
		} else {
			isAuth = true
			authenticator(id)
		}

		return nil, isAuth, nil
	}
}
