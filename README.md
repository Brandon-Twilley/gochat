**GoChat web chat server**

**Authors: Brandon Twilley, Samuel Bretz**

Description: GoChat is a web chat server implemented using golang with a javascript GUI. GoChat works over TCP sockets and HTTP.  Any user can connect by using a standard http connection.

System requirements (Server):

golang >1.7.4

golang websocket (github.com/gorilla/websocket)

golang mux (github.com/gorilla/mux)

**To aquire these:**

`go get github.com/gorilla/mux`

`go get github.com/gorilla/websocket`


*Instructions for hosting web server:*
1. Choose a file directory to host your webserver.
2. Once you’ve picked a directory, upload your client.go, hub.go, main.go, and home.html into your directory and create a directory called files/ but leave empty.  This will be used to host the files uploaded.
3. Once this is uploaded, if you haven’t installed golang, and the additional package, do that already.  Once installed, run ‘go run client.go hub.go main.go’.  Your webserver should now be running on port 8080 by default.
4. Clients can now connect to your server.
5. Find the server’s ip address and have the clients type that into their preferred web browser with port 8080.
6. For example, 192.168.43.203:8080

Client instructions for unicast, blockcast, file upload, changing nicknames and small easter-eggs:

`'message'                     	--Broadcast text 'message’, default case`

`/msg 'username' 'msg'         	--Unicast text 'msg' to 'username'`

`/bmsg 'username' 'msg'        	--Blockcast text, send to all except 'username'`

`/nick 'username'              	--Create/reassign new user.`

`/gucci                        		--Does something interesting`

Uploading a file will result in a new window opening with a URL to the file that was uploaded. The user can choose to send that URL in whatever manner they choose over the chat server.  For example,
`/bmsg <user> http://192.168.43.203:8080/files/abc_123_a.ogg` will send a link to the file to everyone except
`<user>`.