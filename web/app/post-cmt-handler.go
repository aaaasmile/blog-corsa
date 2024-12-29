package app

import (
	"bytes"
	"corsa-blog/idl"
	"html/template"
	"log"
	"net/http"
	"time"
)

func (ph *PostHandler) handleFormNewComment(w http.ResponseWriter, req *http.Request, id string) error {
	lang := req.URL.Query().Get("lang")
	log.Println("process new comment for parent", id, lang)
	// this is coming from a form inside the static page

	// email := req.Body.Get("email")
	// name := req.URL.Query().Get("name")
	// website := req.URL.Query().Get("website")
	templName := "templates/cmt/newcomment.html"
	var partMerged bytes.Buffer
	tmplBody := template.Must(template.New("Body").ParseFiles(templName))
	cmtItem := idl.CmtItem{}
	if err := tmplBody.ExecuteTemplate(&partMerged, "base", cmtItem); err != nil {
		return err
	}

	elapsed := time.Since(ph.start)
	//fmt.Println("response: ", partMerged.String())
	log.Printf("Service total call duration: %v\n", elapsed)
	_, err := w.Write(partMerged.Bytes())
	return err
}
