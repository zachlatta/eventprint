// From github.com/gorilla/websocket
// ---------------------------------
// Copyright 2013 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websockets

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var Hub = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) Broadcast(msg string) {
	h.broadcast <- []byte(msg)
}

func (h *hub) Run() {
	for {
		select {
		case c := <-Hub.register:
			Hub.connections[c] = true
		case c := <-Hub.unregister:
			delete(Hub.connections, c)
			close(c.send)
		case m := <-Hub.broadcast:
			for c := range Hub.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(Hub.connections, c)
				}
			}
		}
	}
}
