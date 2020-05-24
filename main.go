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
	dbMaster := establishDbConnection(config.MysqlDSN)
	dbSlave := establishDbConnection(config.MysqlDSNSlave)
	defer dbMaster.Close()
	defer dbSlave.Close()
	migrateDatabase(dbMaster, "app/migrations")
	templ := registerHtmlTemplates("app/templates")
	sessionWrapper := initSessionStore(config.SessionKey)
	router := initRouter(templ, dbMaster, dbSlave, sessionWrapper)
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

func initRouter(templ *template.Template, dbMaster *sql.DB, dbSlave *sql.DB, sessionWrapper helpers.SessionWrapper) *mux.Router {
	handlerFactory := controller.NewHandlerFactory(templ)

	router := mux.NewRouter()

	router.HandleFunc("/", handlerFactory.CreateHandler(controller.CreateMainGetHandler(templ, dbMaster, dbSlave, sessionWrapper)).ServeHTTP)

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
		handlerFactory.CreateHandler(controller.CreateRegisterPostHandler(templ, dbMaster, sessionWrapper)).ServeHTTP,
	).Methods(http.MethodPost)

	router.HandleFunc(
		"/auth",
		handlerFactory.CreateHandler(controller.CreateAuthGetHandler(templ, sessionWrapper)).ServeHTTP,
	).Methods(http.MethodGet)

	router.HandleFunc(
		"/auth",
		handlerFactory.CreateHandler(controller.CreateAuthPostHandler(templ, dbMaster, sessionWrapper)).ServeHTTP,
	).Methods(http.MethodPost)

	router.HandleFunc(
		"/list",
		handlerFactory.CreateHandler(controller.CreateListGetHandler(templ, dbMaster, sessionWrapper)).ServeHTTP,
	).Methods(http.MethodGet)

	router.HandleFunc(
		"/details",
		handlerFactory.CreateHandler(controller.CreateDetailsGetHandler(templ, dbMaster, sessionWrapper)).ServeHTTP,
	).Methods(http.MethodGet)

	router.HandleFunc(
		"/details",
		handlerFactory.CreateHandler(controller.CreateDetailsPostHandler(templ, dbMaster, sessionWrapper)).ServeHTTP,
	).Methods(http.MethodPost)

	router.HandleFunc(
		"/friends",
		handlerFactory.CreateHandler(controller.CreateFriendsListGetHandler(templ, dbMaster, sessionWrapper)).ServeHTTP,
	).Methods(http.MethodGet)

	router.HandleFunc(
		"/gen",
		handlerFactory.CreateHandler(controller.CreateGeneratorGetHandler(templ, dbMaster, sessionWrapper)).ServeHTTP,
	).Methods(http.MethodGet)

	return router
}

func startWebServer(router *mux.Router, port uint16) {
	err := http.ListenAndServe(":"+strconv.Itoa(int(port)), router)
	if err != nil {
		panic(err.Error())
	}
}
