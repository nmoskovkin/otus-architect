package domain

import "github.com/google/uuid"

type FriendsRepository interface {
	Create(fromUser uuid.UUID, toUser uuid.UUID) error
}
