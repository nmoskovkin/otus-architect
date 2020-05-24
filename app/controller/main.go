package controller

import (
	"architectSocial/app/helpers"
	"architectSocial/app/templates"
	"database/sql"
	"html/template"
	"net/http"
)

func CreateMainGetHandler(templ *template.Template, dbMaster *sql.DB, dbSlave *sql.DB, sessionWrapper helpers.SessionWrapper) ErrorReturningHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		_, id, _ := sessionWrapper.IsAuthenticated(r)
		slave := r.URL.Query().Get("slave")

		var statisticProvider *helpers.UserStatisticProvider
		if slave == "" {
			statisticProvider = helpers.CreateUserStatisticProvider(dbMaster)
		} else {
			statisticProvider = helpers.CreateUserStatisticProvider(dbSlave)
		}

		mostPopularCities, err := statisticProvider.GetMostPopularCities(5)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}
		mostPopularFirstNames, err := statisticProvider.GetMostPopularFirstNames(5)
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
