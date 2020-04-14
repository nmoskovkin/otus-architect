package app

import (
	"architectSocial/app/controller"
	"architectSocial/app/repository"
	"architectSocial/domain"
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"strconv"
)

func NewWebServer(templ *template.Template, port uint16, store sessions.Store, db *sql.DB) error {
	//http.HandleFunc("/register", createRegisterGetHandler(templ))
	//http.Handler("/register")

	s := domain.CreateRegisterUserService(repository.CreateUserRepository(db))
	h1 := controller.CreateRegisterPostHandler(s, templ)
	h2 := controller.CreateRegisterGetHandler(templ)

	router := mux.NewRouter()
	router.HandleFunc("/register", h2.ServeHTTP).Methods("GET")
	router.HandleFunc("/register", h1.ServeHTTP).Methods("POST")

	err := http.ListenAndServe(":"+strconv.Itoa(int(port)), router)
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
