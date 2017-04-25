package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"net"
	"strings"
)

var addr = flag.String("addr", ":8080", "http service address")
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
	http.ServeFile(w, r, "home.html")
}

func main() {

	ipaddress = GetLocalIP();


	gucci = "\n\n########__########_########_########_########_########_\n##_____##_##_______##_______##_______##_______##\n##_____##_##_______##_______##_______##_______##\n########__######___######___######___######___######\n##___##___##_______##_______##_______##_______##\n##____##__##_______##_______##_______##_______##\n##_____##_########_########_########_########_########"

	mod_bot = "HAL_5000-BOT"
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/upload", upload)

	http.HandleFunc("/files/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w,r)
			return			
		}
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
