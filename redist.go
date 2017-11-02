package main

import (
	"fmt"
	"strings"
)

/*
	redist redistributes the messages from client to client.  Here, we take in
	one message, process it to determine who it's for, then we send it off to
	everyone who it belongs to.
*/

func newRedist() *redist {
	return &redist{
		new_client: make(chan *Client),
		broadcast:  make(chan []byte),

		leaving_client: make(chan *Client),
		clients:        make(map[*Client]bool),
	}
}

type redist struct {
	new_client     chan *Client
	clients        map[*Client]bool
	broadcast      chan []byte
	leaving_client chan *Client
}

/*
	redis receives 3 conditions, either a joining client, leaving client, or
	a new message.  if its a new message using chat handles, we parse through
	what we received and distribute the message that way.
*/

func (redis *redist) run() {

	for true {
		select {
		// Unregister requests from clients.

		case cli := <-redis.new_client: //allocates client
			redis.clients[cli] = true
		case cli := <-redis.leaving_client: //deallocates client
			_, notalive := redis.clients[cli]
			if notalive {
				fmt.Println("CLIENT " + cli.name + " LEFT")
				delete(redis.clients, cli)
				close(cli.send)
			}
		case message := <-redis.broadcast: //receives message.  Parses message to send to other clients.
			if is_new_name {
				for cli := range redis.clients {
					select { // links the new client to the chat server
					case cli.send <- message:
					}
				}
			} else if is_unicast {
				for cli := range redis.clients { //SEND ONLY TO THE RECIPIENT AND THE SENDER
					if (strings.Compare(cli.name, recipient) == 0) || (strings.Compare(cli.name, sender) == 0) {
						select {
						case cli.send <- message:
						}
					}
				}
			} else if is_blockcast {
				for cli := range redis.clients { //SEND TO EVERYONE EXCEPT THE RECIPIENT
					if strings.Compare(cli.name, recipient) == 0 {
					} else {
						select {
						case cli.send <- message:
						}
					}
				}

			} else { // TYPICAL CLIENT BROADCAST
				for cli := range redis.clients {
					select {
					case cli.send <- message:
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
