package main

import (
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
)

func main() {
	m := martini.Classic()

	m.Use(martini.Static("public"))
	m.MapTo(Dbm, (*gorp.SqlExecutor)(nil))

	m.Group("/api", func(r martini.Router) {
		r.Get("/hello_world", func() string {
			return "Hello, World!"
		})
	})

	m.Run()
}
