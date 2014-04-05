package main

import (
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/zachlatta/eventprint/server/route"
)

func main() {
	m := martini.Classic()

	m.Use(martini.Static("public"))
	m.MapTo(Dbm, (*gorp.SqlExecutor)(nil))

	m.Group("/api", func(r martini.Router) {
		r.Get("/hello_world", func() string {
			return "Hello, World!"
		})

		r.Group("/attendees", func(r martini.Router) {
			r.Put("/sync", route.Sync)
		})
	})

	m.Run()
}
