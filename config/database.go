package config

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
)

const (
	AppDBName    = "lab_lending"
	DatabasePort = ":5432"
)

func InitPostgresDB() *sql.DB {
	connStr := fmt.Sprintf("postgresql://root@%s?sslmode=disable", DatabasePort)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE DATABASE IF NOT EXISTS ` + AppDBName)
	if err != nil {
		log.Fatal(err)
	}
	db.Close()

	connStr = fmt.Sprintf("postgresql://root@%s/%s?sslmode=disable", DatabasePort, AppDBName)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	ddl, err := ioutil.ReadFile("database/cockroach/ddl.sql")
	if err != nil {
		log.Fatal(err)
	}

	sql := string(ddl)
	_, err = db.Exec(sql)
	if err != nil {
		log.Println(err)
	}

	return db
}
