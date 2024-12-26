package app

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type ServiceHandler struct {
	w     http.ResponseWriter
	req   *http.Request
	debug bool
}

type PostHandler struct {
	debug    bool
	lastPath string
	start    time.Time
}

func (ph *PostHandler) handlePost(w http.ResponseWriter, req *http.Request) error {
	ph.start = time.Now()
	rawbody, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	ph.lastPath = getURLForRoute(req.RequestURI)
	if ph.debug {
		log.Println("[POST] uri requested is: ", ph.lastPath)
		log.Println("[POST] Body is: ", string(rawbody))
	}

	if isPostNewComment(ph.lastPath) {
		return ph.handleFormNewComment(w, req)
	}

	elapsed := time.Since(ph.start)
	log.Printf("Service total call duration: %v\n", elapsed)
	return nil
}

func isPostNewComment(s string) bool {
	return strings.HasSuffix(s, "newcomment")
}

func writeJsonResp(w http.ResponseWriter, resp interface{}) error {
	blobresp, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(blobresp)

	return nil
}
