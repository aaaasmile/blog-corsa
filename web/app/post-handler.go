package app

import (
	"corsa-blog/db"
	"corsa-blog/idl"
	"corsa-blog/web/app/admin"
	"corsa-blog/web/app/comments"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PostHandler struct {
	debug       bool
	lastPath    string
	start       time.Time
	liteDB      *db.LiteDB
	newCmt      chan *idl.CmtItem
	moderateCmt bool
}

func (ph *PostHandler) handlePost(w http.ResponseWriter, req *http.Request) error {
	ph.start = time.Now()
	remPath := ""
	ph.lastPath, remPath = getLastPathInUri(req.RequestURI)
	if ph.debug {
		log.Println("[handlePost] uri requested is: ", ph.lastPath, remPath)
	}

	if parent_id, post_id, ok := isNewComment(ph.lastPath, remPath); ok {
		hc := comments.NewPostCommentHandler(ph.liteDB, ph.debug, ph.moderateCmt, ph.newCmt)
		return hc.HandleFormNewComment(w, req, parent_id, post_id)
	}
	if id, post_id, ok := isDeleteComment(ph.lastPath, remPath); ok {
		hc := comments.NewPostCommentHandler(ph.liteDB, ph.debug, ph.moderateCmt, ph.newCmt)
		return hc.HandleFormDeleteComment(w, req, id, post_id)
	}
	if ok := isAdminReq(ph.lastPath); ok {
		ha := admin.NewAdmin(w, req, ph.liteDB)
		return ha.HandleAdminRequest()
	}

	elapsed := time.Since(ph.start)
	log.Printf("[WARN] ignored request. Total call duration: %v\n", elapsed)
	return nil
}

func isAdminReq(lastPath string) bool {
	if strings.HasPrefix(lastPath, "CallDataService") {
		return true
	}
	return false
}

func isNewComment(lastPath, remPath string) (parent_id int, post_id string, ok bool) {
	// expect something like: /blog-admin/{{.ParentId}}/{{.PostId}}/newcomment?lang=it
	ok = false
	if !strings.HasPrefix(lastPath, "newcomment") {
		return
	}
	arr := strings.Split(remPath, "/")
	if len(arr) > 1 {
		idtxt := arr[len(arr)-2]
		var err error
		if parent_id, err = strconv.Atoi(idtxt); err != nil {
			log.Println("[isNewComment] ERROR parent_id ", err)
			return
		}
		post_id = arr[len(arr)-1]
		ok = true
		return
	}
	return
}

func isDeleteComment(lastPath, remPath string) (id int, post_id string, ok bool) {
	// expect something like: /blog-admin/{{.Id}}/{{.PostId}}/deletecomment?req={{.ReqId}}
	ok = false
	if !strings.HasPrefix(lastPath, "deletecomment") {
		return
	}
	arr := strings.Split(remPath, "/")
	if len(arr) > 1 {
		idtxt := arr[len(arr)-2]
		var err error
		if id, err = strconv.Atoi(idtxt); err != nil {
			log.Println("[isDeleteComment] ERROR id ", err)
			return
		}
		post_id = arr[len(arr)-1]
		ok = true
		return
	}
	return
}

// func writeJsonResp(w http.ResponseWriter, resp interface{}) error {
// 	blobresp, err := json.Marshal(resp)
// 	if err != nil {
// 		return err
// 	}
// 	w.Write(blobresp)

// 	return nil
// }
