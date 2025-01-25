package app

import (
	"corsa-blog/db"
	"corsa-blog/web/app/comments"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type PostHandler struct {
	debug       bool
	lastPath    string
	start       time.Time
	liteDB      *db.LiteDB
	moderateCmt bool
}

func (ph *PostHandler) handlePost(w http.ResponseWriter, req *http.Request) error {
	ph.start = time.Now()
	remPath := ""
	ph.lastPath, remPath = getURLForRoute(req.RequestURI)
	if ph.debug {
		log.Println("[POST] uri requested is: ", ph.lastPath, remPath)
	}

	if id, ok := isPostNewComment(ph.lastPath, remPath); ok {
		hc := comments.NewPostCommentHandler(ph.liteDB, ph.debug, ph.moderateCmt)
		return hc.HandleFormNewComment(w, req, id)
	}

	elapsed := time.Since(ph.start)
	log.Printf("Ignored request. Total call duration: %v\n", elapsed)
	return nil
}

func isPostNewComment(lastPath, remPath string) (string, bool) {
	if !strings.HasPrefix(lastPath, "newcomment") {
		return "", false
	}
	arr := strings.Split(remPath, "/")
	if len(arr) > 0 {
		return arr[len(arr)-1], true
	}
	return "", false
}

func writeJsonResp(w http.ResponseWriter, resp interface{}) error {
	blobresp, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(blobresp)

	return nil
}
