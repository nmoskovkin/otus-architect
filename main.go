package main

import (
	"architectSocial/app"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"html/template"
)

func main() {
	loadEnvFile()
	config := createConfigFromEnvVars()
	db := establishDbConnection(config.MysqlDSN)
	defer db.Close()
	migrateDatabase(db, "app/migrations")
	templ := registerHtmlTemplates("app/templates")
	sessionStore := initSessionStore(config.SessionKey)
	startWebServer(templ, config.Port, sessionStore, db)

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

func initSessionStore(sessionKey string) *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(sessionKey))
}

func startWebServer(tmpl *template.Template, port uint16, store sessions.Store, db *sql.DB) {
	err := app.NewWebServer(tmpl, port, store, db)

	if err != nil {
		panic(err.Error())
	}
}
