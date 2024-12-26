package app

import (
	"log"
	"net/http"
	"time"
)

func (ph *PostHandler) handleFormNewComment(w http.ResponseWriter, req *http.Request) error {
	// this is coming from a form inside the static page
	// lang := req.URL.Query().Get("lang")

	// email := req.Body.Get("email")
	// name := req.URL.Query().Get("name")
	// website := req.URL.Query().Get("website")

	elapsed := time.Since(ph.start)
	log.Printf("Service total call duration: %v\n", elapsed)
	return nil
}
