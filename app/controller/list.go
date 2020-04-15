package controller

import (
	"architectSocial/app/helpers"
	"architectSocial/app/repository"
	"architectSocial/app/templates"
	"database/sql"
	"html/template"
	"net/http"
)

func CreateListGetHandler(templ *template.Template, db *sql.DB, sessionWrapper helpers.SessionWrapper) ErrorReturningHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		repo := repository.CreateMysqlUserRepository(db)
		userList, err := repo.GetAll()
		if err != nil {
			return NewHTTPError(err, 500, "")
		}
		listData := templates.ListData{
			Users:     userList,
			PageTitle: "Users",
		}
		err = templ.ExecuteTemplate(w, "list.html", listData)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}

		return nil
	}
}
