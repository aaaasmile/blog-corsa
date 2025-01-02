package app

import (
	"bytes"
	"corsa-blog/idl"
	"log"
	"net/http"
	"net/mail"
	"text/template"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

func (ph *PostHandler) handleFormNewComment(w http.ResponseWriter, req *http.Request, id string) error {
	lang := req.URL.Query().Get("lang")
	log.Println("process new comment for parent", id, lang)
	err := req.ParseForm()
	if err != nil {
		return err
	}
	email := req.PostFormValue("email")
	name := req.PostFormValue("name")
	commentMd := req.PostFormValue("comment")
	if ph.debug {
		log.Println("orig comment:", commentMd)
		log.Println("name, email:", name, email)
	}
	unsafeComment := blackfriday.Run([]byte(commentMd), blackfriday.WithNoExtensions())
	htmlCmt := bluemonday.UGCPolicy().SanitizeBytes(unsafeComment)
	if ph.debug {
		log.Println("transformed html comment:", string(htmlCmt))
	}

	errMsg := ""
	cmtItem := &idl.CmtItem{
		Email:   email,
		Name:    name,
		Comment: string(htmlCmt),
	}
	if len(htmlCmt) == 0 {
		errMsg = "commento vuoto"
		return ph.renderResNewComment(cmtItem, errMsg, w)
	}

	if email != "" {
		if _, err := mail.ParseAddress(email); err != nil {
			errMsg = "inidirizzo email non valido"
			return ph.renderResNewComment(cmtItem, errMsg, w)
		}
	}
	if name == "" {
		if _, err := mail.ParseAddress(email); err != nil {
			errMsg = "il nome è vuoto"
			return ph.renderResNewComment(cmtItem, errMsg, w)
		}
	}

	return ph.renderResNewComment(cmtItem, errMsg, w)
}

func (ph *PostHandler) renderResNewComment(cmtItem *idl.CmtItem, errMsg string, w http.ResponseWriter) error {
	ctx := struct {
		Cmt       *idl.CmtItem
		ErrMsg    string
		HasErrors bool
	}{
		Cmt:       cmtItem,
		ErrMsg:    errMsg,
		HasErrors: (errMsg != ""),
	}
	//fmt.Println("*** ctx: ", *ctx.Cmt)

	templName := "templates/cmt/newcomment.html"
	var partMerged bytes.Buffer
	tmplBody := template.Must(template.New("Body").ParseFiles(templName))
	if err := tmplBody.ExecuteTemplate(&partMerged, "base", ctx); err != nil {
		return err
	}

	elapsed := time.Since(ph.start)

	log.Printf("Service total call duration: %v\n", elapsed)
	_, err := w.Write(partMerged.Bytes())
	//fmt.Println("response: ", partMerged.String())
	return err
}
