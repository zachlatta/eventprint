package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-martini/martini"
	"github.com/zachlatta/eventprint/server/model"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
)

var Dbm *gorp.DbMap

func init() {
	Dbm = newDbMap()

	Dbm.AddTable(model.Attendee{}).SetKeys(true, "Id")
}

func newDbMap() *gorp.DbMap {
	dialect, driver := dialectAndDriver()
	return &gorp.DbMap{Db: connect(driver), Dialect: dialect}
}

func dialectAndDriver() (gorp.Dialect, string) {
	return gorp.PostgresDialect{}, "postgres"
}

func connect(driver string) *sql.DB {
	dsn := fmt.Sprintf("postgres://docker:docker@%s/docker",
		os.Getenv("DB_1_PORT_5432_TCP_ADDR"))

	if martini.Env == martini.Prod {
		dsn = os.Getenv("DATABASE_URL")
		if dsn == "" {
			panic("DATABASE_URL env variable not set")
		}
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	return db
}
