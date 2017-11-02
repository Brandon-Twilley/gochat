package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var addr = flag.String("addr", ":8000", "http service address")
var ipaddress string

// 10 MB
const MAX_MEMORY = 10 * 1024 * 1024

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusForbidden)
	}

	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()

			path := fmt.Sprintf(fileHeader.Filename)
			fmt.Println(path)
			buf, _ := ioutil.ReadAll(file)
			path = "files/" + path
			ioutil.WriteFile(path, buf, os.ModePerm)

			fmt.Fprintf(w, "URL: %s:8080/%s %s", ipaddress, path, "\n")
		}
	}

}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func serveJS(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	http.ServeFile(w, r, "javascript/jscript.js")
}

func serveStyle(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	http.ServeFile(w, r, "style/style.css")
}

func main() {

	ipaddress = GetLocalIP()

	gucci = "\n\n########__########_########_########_########_########_\n##_____##_##_______##_______##_______##_______##\n##_____##_##_______##_______##_______##_______##\n########__######___######___######___######___######\n##___##___##_______##_______##_______##_______##\n##____##__##_______##_______##_______##_______##\n##_____##_########_########_########_########_########"

	flag.Parse()
	red := newRedist()
	go red.run()
	/*  "/upload" tells us where our file was uploaded.  This URL can only be viewed
	 	by the person that uploaded the document.  From there, they can choose to
		forward the URL to any other user using the client.
	*/
	http.HandleFunc("/upload", upload)

	/*
		/files/ holds the documents that have been uploaded by the user.  This is held
		in a subdirectory of the root server by the same name.  If anyone tries to access
		localhost:8080/files, they are returned with a 404 error.  This prevents people
		from looking at all the contents stored in the /file/ subdirectory of our webserver.
	*/
	http.HandleFunc("/files/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	/*
		This is the root directory of the main webserver.  This loads up the homepage for our
		client.  From here, the /ws deals with communications between our sockets on our client
		javascript functions and the sockets on our server.  The ws handlefunc deals with most of
		the networking in the system (in the serveWs function).
	*/
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/javascript/jscript.js", serveJS)
	http.HandleFunc("/style/style.css", serveStyle)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websox(red, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
