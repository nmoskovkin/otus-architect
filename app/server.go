package app

import (
	"architectSocial/app/controller"
	"database/sql"
	"errors"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"strconv"
)

func NewWebServer(templ *template.Template, port uint16, store sessions.Store, db *sql.DB) error {
	//http.HandleFunc("/register", createRegisterGetHandler(templ))
	//http.Handler("/register")

	h := controller.CreateRegisterHandler(templ, db)
	http.HandleFunc("/register", h.ServeHTTP)

	err := http.ListenAndServe(":"+strconv.Itoa(int(port)), nil)
	if err != nil {
		return errors.New("Failed to create a web server. Error: " + err.Error())
	}

	return nil
}

//func showErrorPage(w http.ResponseWriter, templ *template.Template) {
//	_ = templ.ExecuteTemplate(w, "error.html", struct{ PageTitle string }{
//		PageTitle: "Register New User",
//	})
//}
