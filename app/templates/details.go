package templates

import "architectSocial/app/repository"

type DetailsData struct {
	PageTitle       string
	User            repository.UserItem
	IsAlreadyFriend bool
	Errors          []string
	CurrentUserId   string
}
