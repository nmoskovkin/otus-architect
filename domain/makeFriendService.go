package domain

import (
	"architectSocial/domain/helper"
	"fmt"
	"github.com/google/uuid"
)

type MakeFriendDto struct {
	From string
	To   string
}

type MakeFriendService func(dto *MakeFriendDto) (*helper.ValidationResult, error)

func makeFriendValidateDto(dto *MakeFriendDto) *helper.ValidationResult {
	result := helper.NewValidationResult()
	if dto.From == "" {
		result.AddError("From", "From is empty")
	}
	if dto.To == "" {
		result.AddError("To", "To is empty")
	}
	if dto.From != "" && dto.To != "" && dto.From == dto.To {
		// TODO: Tricks with uuid?
		result.AddError("To", "To and from are same")
	}

	if dto.From != "" {
		_, err := uuid.Parse(dto.From)
		if err != nil {
			result.AddError("From", "From must be correct uuid")
		}
	}
	if dto.To != "" {
		_, err := uuid.Parse(dto.From)
		if err != nil {
			result.AddError("From", "From must be correct uuid")
		}
	}

	return result
}

func CreateMakeFriendService(friendsRepository FriendsRepository) MakeFriendService {
	return func(dto *MakeFriendDto) (*helper.ValidationResult, error) {
		validationResult := makeFriendValidateDto(dto)
		if !validationResult.IsValid() {
			return validationResult, nil
		}

		uuidFrom, err := uuid.Parse(dto.From)
		if err != nil {
			return nil, fmt.Errorf("failed to parse field \"From\", error: %s", err.Error())
		}
		uuidTo, err := uuid.Parse(dto.To)
		if err != nil {
			return nil, fmt.Errorf("failed to parse field \"To\", error: %s", err.Error())
		}
		err = friendsRepository.Create(uuidFrom, uuidTo)
		if err != nil {
			return nil, fmt.Errorf("failed to create friendship: %s", err.Error())
		}

		return nil, nil
	}
}
