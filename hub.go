// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main
import (
	"strings"
)
// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

var mod_bot string

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}


func (h *Hub) run() {

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:

			if is_new_name {

				for client := range h.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}

			} else if (is_unicast) {
				for client := range h.clients {	//SEND ONLY TO THE RECIPIENT AND THE SENDER
					if ((strings.Compare(client.name,recipient) == 0) || (strings.Compare(client.name,sender) == 0)) {
						select {
						case client.send <- message:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
				}

			} else if (is_blockcast) {
				for client := range h.clients { //SEND TO EVERYONE EXCEPT THE RECIPIENT
					if ((strings.Compare(client.name, recipient) == 0)) {

					} else {
						select {
						case client.send <- message:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}

				}

			} else {
				for client := range h.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
		is_gucci = false
		is_new_name = false
		is_blockcast = false
		is_unicast = false

		sender = ""
		recipient = ""
	}
}