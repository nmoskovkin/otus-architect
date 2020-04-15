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

	//fmt.Println(err)
	//fmt.Println(config)

	//db, err := sql.Open("mysql", config.MysqlDSN)
	//r, err1 := db.Query("SELECT 1,2")
	//fmt.Println(r)
	//fmt.Println(err)
	//fmt.Println(err1)
	//
	//var col1, col2 []byte
	//for r.Next() {
	//	// Scan the value to []byte
	//	err = r.Scan(&col1, &col2)
	//
	//	if err != nil {
	//		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	//	}
	//
	//	// Use the string value
	//	fmt.Println(string(col1), string(col2))
	//}

	//t, _ := template.New("header").ParseFiles("templates/header.html")
	//t.New("footer").ParseFiles("templates/footer.html")
	//t.New("list").ParseFiles("templates/list.html")

	//tmpl, err := template.ParseFiles("templates/list.html")

	//driver, err := postgres.WithInstance(db, &postgres.Config{})
	//m, err := migrate.NewWithDatabaseInstance(
	//	"file:///migrations",
	//	"postgres", driver)
	//m.Steps(2)

	//if value, exists := os.LookupEnv(key); exists {
	//	return value
	//}
	//
	//const name, age = "Kim", 22
	//fmt.Printf("%s is %d years old.\n", name, age)
	//
	//var cfg Config
	//err := envconfig.Process("", &cfg)
	//if err != nil {
	//
	//	//processError(err)
	//}
	//
	//// It is conventional not to worry about any
	//// error returned by Printf.
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
	router.HandleFunc(
		"/register",
		handlerFactory.CreateHandler(controller.CreateRegisterGetHandler(templ)).ServeHTTP,
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

	return router
}

func startWebServer(router *mux.Router, port uint16) {
	err := http.ListenAndServe(":"+strconv.Itoa(int(port)), router)
	if err != nil {
		panic(err.Error())
	}
}
