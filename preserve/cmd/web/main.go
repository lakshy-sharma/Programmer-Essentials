package main

import (
	"database/sql"
	"flag"
	"net/http"
	"preserve/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type Preserve struct {
	logger logrus.Logger
	notes  *mysql.NoteModel
}

type Config struct {
	Port      string
	StaticDir string
}

func createLogger() *logrus.Logger {
	preserveLogger := logrus.New()
	return preserveLogger
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	configs := new(Config)

	// Read command line options while starting the application.
	flag.StringVar(&configs.Port, "port", ":4000", "Network Port To Host the Application")
	flag.StringVar(&configs.StaticDir, "static-dir", "./ui/static", "Path to the Application's Static assets.")
	dsn := flag.String("dsn", "web:preserve123@tcp(127.0.0.1:3306)/preserve?parseTime=true", "MySQL data connection string")

	flag.Parse()

	db, err := openDB(*dsn)
	// Initialize the loggers and notes and add them to the application.
	preserve := &Preserve{
		logger: *createLogger(),
		notes:  &mysql.NoteModel{DB: db},
	}
	if err != nil {
		preserve.logger.Fatal(err)
	}
	preserve.logger.Println("Setup connection with DB successfully.")
	defer db.Close()

	// Creating a http server with dedicated routes.
	srv := &http.Server{
		Addr:    configs.Port,
		Handler: preserve.routes(),
	}

	// Starting the Application Server.
	preserve.logger.Infof("Starting Application on Port %s", configs.Port)
	err = srv.ListenAndServe()
	preserve.logger.Fatal(err)
}
