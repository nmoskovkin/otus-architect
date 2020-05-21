package controller

import (
	"architectSocial/app/helpers"
	"architectSocial/app/templates"
	"database/sql"
	"html/template"
	"net/http"
)

func CreateMainGetHandler(templ *template.Template, db *sql.DB, sessionWrapper helpers.SessionWrapper) ErrorReturningHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		_, id, _ := sessionWrapper.IsAuthenticated(r)
		statisticProvier := helpers.CreateUserStatisticProvider(db)
		mostPopularCities, err := statisticProvier.GetMostPopularCities(5)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}
		mostPopularFirstNames, err := statisticProvier.GetMostPopularFirstNames(5)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}

		mainData := templates.MainData{
			PageTitle:             "Main Page",
			CurrentUserId:         id,
			MostPopularCities:     mostPopularCities,
			MostPopularFirstNames: mostPopularFirstNames,
		}
		err = templ.ExecuteTemplate(w, "main.html", mainData)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}

		return nil
	}
}
