package main

import (
	"log"
	"net/http"
	"time"

	"encoding/json"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	r *redist

	// The websocket connection.
	conn *websocket.Conn
	name string
	// Buffered channel of outbound messages.
	send chan []byte
}

// reader pumps messages from the websocket connection to the hub.
//
// The application runs reader in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.

type msg struct {
	User string `json:"user"`
	Text string `json:"text"`
}

var is_unicast bool
var is_blockcast bool
var is_new_name bool
var is_gucci bool

var sender string
var recipient string
var gucci string

func irccall(m *msg, c *Client) {
	/*
		Expected data formats
		'message'                         --Broadcast text 'message, default case
		/msg 'username' 'msg'             --Unicast text 'message' to 'username'
		/bmsg 'username' 'msg'            --Blockcast text, send to all except 'username'
		/bcfile 'filepath'                --Broadcast file to all users, where 'filepath' is
				full path on localmachine
		/ucfile 'username' 'filepath'     --Unicast file to 'username', where 'filepath' is
				full path on localmachine
		/nick 'username'                  --Create/reassign new user.
		/gucci                            --Does something interesting*/

	partition := strings.Fields(m.Text)
	fmt.Println("PARTITION: ", partition)
	for i := 0; i < len(partition); i++ {
		fmt.Println(i, " : ", partition[i])
	}
	if strings.Compare(partition[0], "/msg") == 0 {
		is_unicast = true
		is_blockcast = false
		is_new_name = false
		is_gucci = false

		sender = m.User
		if len(partition) < 2 {
			return
		}
		recipient = partition[1]

	} else if (strings.Compare(partition[0], "/bmsg")) == 0 {
		is_unicast = false
		is_blockcast = true
		is_new_name = false
		is_gucci = false
		if len(partition) < 2 {
			return
		}
		sender = m.User
		recipient = partition[1]
	} else if (strings.Compare(partition[0], "/nick")) == 0 {
		is_unicast = false
		is_blockcast = false
		is_new_name = true
		is_gucci = false

		sender = m.User
		if len(partition) < 2 {
			return
		}
		m.User = partition[1]
		c.name = partition[1]
		fmt.Println("NEW NICKNAME: ", m.User)
	} else if (strings.Compare(partition[0], "/gucci")) == 0 {
		is_unicast = false
		is_blockcast = false
		is_new_name = false
		is_gucci = true

		m.Text = gucci
	}

	return
}

// 	Each instance of this function is running for all clients connected
// 	to the webserver.
func (cli *Client) reader() {
	defer func() {
		cli.r.leaving_client <- cli
		cli.conn.Close()
	}()
	cli.conn.SetReadLimit(maxMessageSize)
	for {
		_, message, err := cli.conn.ReadMessage()
		m := &msg{User: cli.name, Text: string(message)}
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		if message[0] == '/' {
			irccall(m, cli)
		}
		str_message, err := json.Marshal(m)

		message = str_message
		fmt.Println(string(message))
		//send message read to redist
		cli.r.broadcast <- message
	}
}

// 	Each instance of this function is running for all clients connected
// 	to the webserver.
func (cli *Client) writer() {
	//retrieve message read from redist and determine if client is connected
	for {
		select {
		case message, ok := <-cli.send:
			cli.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Channel closed.
				cli.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := cli.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			// Add queued chat messages to the current websocket message.
			n := len(cli.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-cli.send)
			}
			if err := w.Close(); err != nil {
				return
			}
			//constantly pings our client to see if there still exists a connection.
			//if there exists no connection, terminate.
		}
	}
}

// websox handles websocket requests from the peer.
func websox(red *redist, w http.ResponseWriter, r *http.Request) {

	// this gives back the client a specific socket to run on.
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}
	//	we create the client datatype to keep track of all our clients in
	//	the system and store it in a map.
	c := &Client{r: red, conn: conn, send: make(chan []byte, 256)}

	//	the initial name of our client is their IP address and socket.
	//	They can change this if they use the /nick command.
	c.name = c.conn.RemoteAddr().String()

	red.new_client <- c
	go c.writer()
	c.reader()
}
