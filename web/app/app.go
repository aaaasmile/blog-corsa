package app

import (
	"corsa-blog/conf"
	"fmt"
	"log"
	"net/http"
)

func APiHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		status := http.StatusOK
		gh := GetHandler{
			debug: conf.Current.Debug,
		}
		if err := gh.handleGet(w, req, &status); err != nil {
			log.Println("Error on process request: ", err)
			if status == http.StatusNotFound {
				http.Error(w, "404 - Not found", http.StatusNotFound)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}
	case "POST":
		ph := PostHandler{
			debug: conf.Current.Debug,
		}
		if err := ph.handlePost(w, req); err != nil {
			log.Println("[POST] Error: ", err)
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
	}
}
