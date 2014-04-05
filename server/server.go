package main

import "github.com/go-martini/martini"

func main() {
	m := martini.Classic()

	m.Use(martini.Static("public"))

	m.Group("/api", func(r martini.Router) {
		r.Get("/hello_world", func() string {
			return "Hello, World!"
		})
	})

	m.Run()
}
