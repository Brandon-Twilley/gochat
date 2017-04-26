
package main
import (
	"strings"
)

/*
	redist redistributes the messages from client to client.  Here, we take in
	one message, process it to determine who it's for, then we send it off to
	everyone who it belongs to.
*/


func newRedist() *redist {
	return &redist{
		new_client:   make(chan *Client),
		broadcast:  make(chan []byte),

		leaving_client: make(chan *Client),
		clients:    make(map[*Client]bool),


	}
}

type redist struct {

	// Register requests from the clients.
	new_client chan *Client

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Unregister requests from clients.
	leaving_client chan *Client

}

func (redis *redist) run() {

	for {
		select {
		case cli := <-redis.new_client:
			redis.clients[cli] = true
		case cli := <-redis.leaving_client:
			if _, ok := redis.clients[cli]; ok {
				delete(redis.clients, cli)
				close(cli.send)
			}
		case message := <-redis.broadcast:

			if is_new_name {

				for cli := range redis.clients {
					select {	// links the new client to the chat server
					case cli.send <- message:
					default:
						close(cli.send)
						delete(redis.clients, cli)
					}
				}

			} else if (is_unicast) {
				for cli := range redis.clients {	//SEND ONLY TO THE RECIPIENT AND THE SENDER
					if ((strings.Compare(cli.name,recipient) == 0) || (strings.Compare(cli.name,sender) == 0)) {
						select {
						case cli.send <- message:
						default:
							close(cli.send)
							delete(redis.clients, cli)
						}
					}
				}

			} else if (is_blockcast) {
				for cli := range redis.clients { //SEND TO EVERYONE EXCEPT THE RECIPIENT
					if ((strings.Compare(cli.name, recipient) == 0)) {

					} else {
						select {
						case cli.send <- message:
						default:
							close(cli.send)
							delete(redis.clients, cli)
						}
					}

				}

			} else {		// TYPICAL CLIENT BROADCAST
				for cli := range redis.clients {
					select {
					case cli.send <- message:
					default:
						close(cli.send)
						delete(redis.clients, cli)
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