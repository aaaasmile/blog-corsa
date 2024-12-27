package app

import (
	"corsa-blog/conf"
	"corsa-blog/idl"
	"corsa-blog/util"
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
	debug    bool
	lastPath string
}

func (gh *GetHandler) handleGet(w http.ResponseWriter, req *http.Request, status *int) error {
	u, _ := url.Parse(req.RequestURI)

	log.Println("GET requested ", u)

	gh.lastPath = getURLForRoute(req.RequestURI)
	if gh.debug {
		log.Println("Check the last path ", gh.lastPath)
	}

	if isRootPattern(gh.lastPath) {
		return gh.handleGetApp(w)
	}
	if isValidateEmail(gh.lastPath) {
		return gh.handleGetValidateEmail(w, req)
	}

	*status = http.StatusNotFound
	return fmt.Errorf("invalid GET request for %s", gh.lastPath)

}

func isValidateEmail(lastPath string) bool {
	return strings.HasPrefix(lastPath, "validatoremail")
}

func isRootPattern(lastPath string) bool {
	str := strings.ReplaceAll(conf.Current.RootURLPattern, "/", "")
	return strings.HasPrefix(lastPath, str)
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

	err := tmplIndex.ExecuteTemplate(w, "base", pagectx)
	if err != nil {
		return err
	}
	return nil
}

func getURLForRoute(uri string) string {
	arr := strings.Split(uri, "/")
	//fmt.Println("split: ", arr, len(arr))
	for i := len(arr) - 1; i >= 0; i-- {
		ss := arr[i]
		if ss != "" {
			if !strings.HasPrefix(ss, "?") {
				//fmt.Printf("Url for route is %s\n", ss)
				return ss
			}
		}
	}
	return uri
}
