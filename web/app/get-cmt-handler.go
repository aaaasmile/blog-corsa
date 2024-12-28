package app

import (
	"bytes"
	"corsa-blog/idl"
	"html/template"
	"log"
	"net/http"
)

func (gh *GetHandler) handleComments(w http.ResponseWriter, req *http.Request, id string) error {
	lang := req.URL.Query().Get("lang")
	log.Printf("get comments for id=%s, lang=%s", id, lang)

	// TODO read comments from data file

	templName := "templates/get/comments.html"
	var partHeader, partTree, partFoot, partMerged bytes.Buffer
	tmplBody := template.Must(template.New("DocPart").ParseFiles(templName))
	cmtItem := idl.CmtItem{}

	if err := tmplBody.ExecuteTemplate(&partHeader, "head", cmtItem); err != nil {
		return err
	}
	if err := tmplBody.ExecuteTemplate(&partTree, "tree", cmtItem); err != nil {
		return err
	}

	if err := tmplBody.ExecuteTemplate(&partFoot, "foot", cmtItem); err != nil {
		return err
	}
	partHeader.WriteTo(&partMerged)
	partTree.WriteTo(&partMerged)
	partFoot.WriteTo(&partMerged)

	_, err := w.Write(partMerged.Bytes())

	return err
}
