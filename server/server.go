package main

import (
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/zachlatta/eventprint/server/route"
	"github.com/zachlatta/eventprint/server/websockets"
)

func main() {
	go websockets.Hub.Run()

	m := martini.Classic()

	m.Use(martini.Static("public"))
	m.MapTo(Dbm, (*gorp.SqlExecutor)(nil))

	m.Group("/api", func(r martini.Router) {
		r.Get("/ws", websockets.ServeWs)

		r.Group("/attendees", func(r martini.Router) {
			r.Put("/sync", route.Sync)

			r.Get("", route.GetAttendees)
			r.Group("/:id", func(r martini.Router) {
				r.Get("", route.GetAttendee)
				r.Put("/check_in", route.CheckInAttendee)
			})
		})
	})

	m.Run()
}
