package templates

import "architectSocial/app/helpers"

type MainData struct {
	PageTitle             string
	CurrentUserId         string
	MostPopularCities     []helpers.PopularCity
	MostPopularFirstNames []helpers.PopularFirstName
}
