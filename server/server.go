package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/zachlatta/eventprint/server/eventbrite"
)

func main() {
	m := martini.Classic()

	m.Use(martini.Static("public"))
	m.MapTo(Dbm, (*gorp.SqlExecutor)(nil))

	m.Group("/api", func(r martini.Router) {
		r.Get("/hello_world", func() string {
			return "Hello, World!"
		})

		r.Get("/eventbrite", func(w http.ResponseWriter) string {
			attendees, err := eventbrite.GetAttendees()
			if err != nil {
				log.Fatal(err)
			}

			json, err := json.Marshal(attendees)
			if err != nil {
				log.Fatal(err)
			}

			w.Header().Set("Content-Type", "application/json")
			return string(json)
		})
	})

	m.Run()
}
