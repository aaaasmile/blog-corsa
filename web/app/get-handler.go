package app

import (
	"corsa-blog/conf"
	"corsa-blog/db"
	"corsa-blog/idl"
	"corsa-blog/util"
	"corsa-blog/web/app/comments"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
)

type PageCtx struct {
	RootUrl        string
	Buildnr        string
	ServerName     string
	VuetifyLibName string
	VueLibName     string
}

type GetHandler struct {
	debug         bool
	lastPath      string
	liteCommentDB *db.LiteDB
}

func (gh *GetHandler) handleGet(w http.ResponseWriter, req *http.Request, status *int) error {
	u, _ := url.Parse(req.RequestURI)

	log.Println("GET requested ", u)

	remPath := ""
	gh.lastPath, remPath = getLastPathInUri(req.RequestURI)
	if gh.debug {
		log.Println("Check the last path ", gh.lastPath, remPath)
	}

	if isHeadNav(gh.lastPath) {
		return gh.handleHeadNav(w)
	}

	if isRootPattern(gh.lastPath) {
		return gh.handleGetApp(w)
	}
	if isValidateEmail(gh.lastPath) {
		return gh.handleGetValidateEmail(w, req)
	}
	if post_id, ok := isComments(gh.lastPath, remPath); ok {
		hc := comments.NewGetCommentHandler(gh.liteCommentDB, gh.debug)
		return hc.HandleCommentsTitle(w, req, post_id)
	}
	if post_id, ok := isCommentDetails(gh.lastPath, remPath); ok {
		hc := comments.NewGetCommentHandler(gh.liteCommentDB, gh.debug)
		return hc.HandleCommentsDetails(w, req, post_id)
	}
	if id, ok := isFormForReplyComment(gh.lastPath, remPath); ok {
		hc := comments.NewGetCommentHandler(gh.liteCommentDB, gh.debug)
		return hc.HandleFormForReplyComment(w, req, id)
	}

	*status = http.StatusNotFound
	return fmt.Errorf("[WARN] invalid GET request for %s", gh.lastPath)

}

func isValidateEmail(lastPath string) bool {
	return strings.HasPrefix(lastPath, "validatoremail")
}

func isComments(lastPath string, remPath string) (string, bool) {
	if !strings.HasPrefix(lastPath, "comments") {
		return "", false
	}
	arr := strings.Split(remPath, "/")
	if len(arr) > 0 {
		return arr[len(arr)-1], true
	}
	return "", false
}

func isCommentDetails(lastPath string, remPath string) (string, bool) {
	if !strings.HasPrefix(lastPath, "cmtDetails") {
		return "", false
	}
	arr := strings.Split(remPath, "/")
	if len(arr) > 0 {
		return arr[len(arr)-1], true
	}
	return "", false
}

func isFormForReplyComment(lastPath string, remPath string) (string, bool) {
	if !strings.HasPrefix(lastPath, "cmtform") {
		return "", false
	}
	arr := strings.Split(remPath, "/")
	if len(arr) > 0 {
		return arr[len(arr)-1], true
	}
	return "", false
}

func isRootPattern(lastPath string) bool {
	str := strings.ReplaceAll(conf.Current.RootURLPattern, "/", "")
	return strings.HasPrefix(lastPath, str)
}

func isHeadNav(lastPath string) bool {
	return strings.HasPrefix(lastPath, "headnav")
}

func (gh *GetHandler) handleGetValidateEmail(w http.ResponseWriter, req *http.Request) error {
	email := req.URL.Query().Get("email")
	lang := req.URL.Query().Get("lang")
	if gh.debug {
		log.Printf("email to validate is %s, language: %s", email, lang)
	}
	if email == "" {
		return nil
	}
	if _, err := mail.ParseAddress(email); err != nil {
		if gh.debug {
			log.Println("email is invalid", err)
		}
		valid_err := "Indirizzo Email non valido"
		w.Write([]byte(valid_err))
	}

	return nil
}

func (gh *GetHandler) handleGetApp(w http.ResponseWriter) error {
	if gh.debug {
		log.Println("provides the App dashboard")
	}
	w.Header().Set("Cache-Control", "stale-while-revalidate=3600")
	pagectx := PageCtx{
		RootUrl:        conf.Current.RootURLPattern,
		Buildnr:        idl.Buildnr,
		ServerName:     conf.Current.ServerName,
		VuetifyLibName: conf.Current.VuetifyLibName,
		VueLibName:     conf.Current.VueLibName,
	}

	templName := "templates/vue/index.html"

	tmplIndex := template.Must(template.New("AppIndex").ParseFiles(util.GetFullPath(templName)))

	return tmplIndex.ExecuteTemplate(w, "base", pagectx)
}

func (gh *GetHandler) handleHeadNav(w http.ResponseWriter) error {
	if gh.debug {
		log.Println("provides the header nav")
	}
	templName := "templates/get/headnav.html"
	tmplIndex := template.Must(template.New("Nav").ParseFiles(util.GetFullPath(templName)))
	return tmplIndex.ExecuteTemplate(w, "headnav", struct{}{})
}

func getLastPathInUri(uri string) (string, string) {
	arr := strings.Split(uri, "/")
	for i := len(arr) - 1; i >= 0; i-- {
		last := arr[i]
		rem_ix := i
		if last != "" {
			if !strings.HasPrefix(last, "?") {
				remPath := strings.Join(arr[0:rem_ix], "/")
				return last, remPath
			}
		}
	}
	return uri, ""
}
