package main

import (
	"architectSocial/app/controller"
	"architectSocial/app/helpers"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"html/template"
	"net/http"
	"strconv"
)

func main() {
	loadEnvFile()
	config := createConfigFromEnvVars()
	db := establishDbConnection(config.MysqlDSN)
	defer db.Close()
	migrateDatabase(db, "app/migrations")
	templ := registerHtmlTemplates("app/templates")
	sessionWrapper := initSessionStore(config.SessionKey)
	router := initRouter(templ, db, sessionWrapper)
	startWebServer(router, config.Port)
}

func loadEnvFile() {
	_ = godotenv.Load()
}

func createConfigFromEnvVars() *Config {
	config, err := NewConfig()
	if err != nil {
		panic(err.Error())
	}

	return config
}

func establishDbConnection(mysqlDsn string) *sql.DB {
	db, err := sql.Open("mysql", mysqlDsn)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	return db
}

func migrateDatabase(db *sql.DB, path string) {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		panic(err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+path,
		"mysql", driver)

	if err != nil {
		panic(err.Error())
	}

	err = m.Up()
	if err != nil && err.Error() != migrate.ErrNoChange.Error() {
		panic(err.Error())
	}
}

func registerHtmlTemplates(path string) *template.Template {
	tmpl, err := template.ParseGlob("./" + path + "/*")

	if err != nil {
		panic(err.Error())
	}

	return tmpl
}

func initSessionStore(sessionKey string) helpers.SessionWrapper {
	store := sessions.NewCookieStore([]byte(sessionKey))

	return helpers.NewGorillaSessionWrapper(store)
}

func initRouter(templ *template.Template, db *sql.DB, sessionWrapper helpers.SessionWrapper) *mux.Router {
	handlerFactory := controller.NewHandlerFactory(templ)

	router := mux.NewRouter()

	router.HandleFunc("/", handlerFactory.CreateHandler(controller.CreateMainGetHandler(templ, db, sessionWrapper)).ServeHTTP)

	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		_ = sessionWrapper.Clear(r, w)
		http.Redirect(w, r, "/list", 302)
	}).Methods(http.MethodPost)

	router.HandleFunc(
		"/register",
		handlerFactory.CreateHandler(controller.CreateRegisterGetHandler(templ, sessionWrapper)).ServeHTTP,
	).Methods(http.MethodGet)

	router.HandleFunc(
		"/register",
		handlerFactory.CreateHandler(controller.CreateRegisterPostHandler(templ, db, sessionWrapper)).ServeHTTP,
	).Methods(http.MethodPost)

	router.HandleFunc(
		"/auth",
		handlerFactory.CreateHandler(controller.CreateAuthGetHandler(templ, sessionWrapper)).ServeHTTP,
	).Methods(http.MethodGet)

	router.HandleFunc(
		"/auth",
		handlerFactory.CreateHandler(controller.CreateAuthPostHandler(templ, db, sessionWrapper)).ServeHTTP,
	).Methods(http.MethodPost)

	router.HandleFunc(
		"/list",
		handlerFactory.CreateHandler(controller.CreateListGetHandler(templ, db, sessionWrapper)).ServeHTTP,
	).Methods(http.MethodGet)

	router.HandleFunc(
		"/details",
		handlerFactory.CreateHandler(controller.CreateDetailsGetHandler(templ, db, sessionWrapper)).ServeHTTP,
	).Methods(http.MethodGet)

	router.HandleFunc(
		"/details",
		handlerFactory.CreateHandler(controller.CreateDetailsPostHandler(templ, db, sessionWrapper)).ServeHTTP,
	).Methods(http.MethodPost)

	router.HandleFunc(
		"/friends",
		handlerFactory.CreateHandler(controller.CreateFriendsListGetHandler(templ, db, sessionWrapper)).ServeHTTP,
	).Methods(http.MethodGet)

	router.HandleFunc(
		"/gen",
		handlerFactory.CreateHandler(controller.CreateGeneratorGetHandler(templ, db, sessionWrapper)).ServeHTTP,
	).Methods(http.MethodGet)

	return router
}

func startWebServer(router *mux.Router, port uint16) {
	err := http.ListenAndServe(":"+strconv.Itoa(int(port)), router)
	if err != nil {
		panic(err.Error())
	}
}
