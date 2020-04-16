package templates

import "architectSocial/app/repository"

type ListData struct {
	PageTitle     string
	Users         []repository.FindAllItem
	CurrentUserId string
}
