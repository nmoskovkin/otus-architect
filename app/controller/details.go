package controller

import (
	"architectSocial/app/helpers"
	"architectSocial/app/repository"
	"architectSocial/app/templates"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

func CreateDetailsGetHandler(templ *template.Template, db *sql.DB, sessionWrapper helpers.SessionWrapper) ErrorReturningHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := r.URL.Query().Get("id")
		if id == "" {
			return NewHTTPError(fmt.Errorf("user not found"), 404, "")
		}

		repo := repository.CreateMysqlUserRepository(db)
		userList, err := repo.GetAll(repository.GetAllFilter{Id: id})
		if len(userList) == 0 {
			return NewHTTPError(fmt.Errorf("user not found"), 404, "")
		}

		if err != nil {
			return NewHTTPError(err, 500, "")
		}
		detailsData := templates.DetailsData{
			User:      userList[0],
			PageTitle: "Users",
		}
		err = templ.ExecuteTemplate(w, "details.html", detailsData)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}

		return nil
	}
}
