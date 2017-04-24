// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dlion/goimgur"
)

var addr = flag.String("addr", ":8080", "http service address")
var ipaddress string

// 1MB
const MAX_MEMORY = 1 * 1024 * 1024

type img_info struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type resp struct {
	Data *img_info `json:"data"`
}

func upload_to_imgur(filepath string) (url_path string) {
	str, err := goImgur.Upload(filepath, "8f4bb2ff3d41947")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(*str)

	byte_resp := []byte(*str)
	r := &resp{}
	json.Unmarshal(byte_resp, r)
	url_path = "https://www.imgur.com/" + r.Data.Id
	return url_path

}

func upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusForbidden)
	}

	url := ""

	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()

			/*
				TODO: IMPORTANT!!! Change this directory below when you start a server.
				This should point to your home directory with the subdirectory of your
				file.
			*/

			path := fmt.Sprintf(fileHeader.Filename)
			fmt.Println(path)
			buf, _ := ioutil.ReadAll(file)
			path = "files/" + path
			ioutil.WriteFile(path, buf, os.ModePerm)
			url = upload_to_imgur(path)

			fmt.Fprintf(w, "URL: %s/%s %s %s %s", ipaddress, path, "\n", "IMGUR: ", url)
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

	ipaddress = "192.168.43.203:8080"
	gucci = "\n\n########__########_########_########_########_########_\n##_____##_##_______##_______##_______##_______##\n##_____##_##_______##_______##_______##_______##\n########__######___######___######___######___######\n##___##___##_______##_______##_______##_______##\n##____##__##_______##_______##_______##_______##\n##_____##_########_########_########_########_########"

	mod_bot = "HAL_5000-BOT"
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/upload", upload)

	http.HandleFunc("/files/", func(w http.ResponseWriter, r *http.Request) {
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
