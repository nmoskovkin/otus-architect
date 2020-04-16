package templates

import "architectSocial/app/repository"

type DetailsData struct {
	PageTitle       string
	User            repository.FindAllItem
	IsAlreadyFriend bool
	Errors          []string
	CurrentUserId   string
}
