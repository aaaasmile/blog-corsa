package comments

import (
	"bytes"
	"corsa-blog/db"
	"corsa-blog/idl"
	"log"
	"net/http"
	"text/template"
	"time"
)

func NewGetCommentHandler(liteDB *db.LiteDB, debug bool) *CommentHandler {
	res := &CommentHandler{
		debug:  debug,
		liteDB: liteDB,
		start:  time.Now(),
	}
	return res
}

func (ch *CommentHandler) HandleComments(w http.ResponseWriter, req *http.Request, post_id string) error {
	lang := req.URL.Query().Get("lang")
	log.Printf("get comments for id=%s, lang=%s", post_id, lang)

	cmtNode, err := ch.liteDB.GeCommentsForPostId(post_id)
	if err != nil {
		return nil
	}

	templName := "templates/cmt/get-comments.html"
	var partHeader, partTree, partFoot, partMerged bytes.Buffer
	tmplBody := template.Must(template.New("DocPart").ParseFiles(templName))

	ctxHead := struct {
		CmtTotText string
	}{
		CmtTotText: cmtNode.GetTextNumComments(),
	}
	if err := tmplBody.ExecuteTemplate(&partHeader, "head", ctxHead); err != nil {
		return err
	}

	ctxTree := struct {
		CmtLines []string
	}{
		CmtLines: []string{"<li>Risposta 1.2 <button>Rispondi</button></li>"},
	}
	if err := tmplBody.ExecuteTemplate(&partTree, "tree", ctxTree); err != nil {
		return err
	}

	cmtItem := idl.CmtItem{}
	if err := tmplBody.ExecuteTemplate(&partFoot, "foot", cmtItem); err != nil {
		return err
	}
	partHeader.WriteTo(&partMerged)
	partTree.WriteTo(&partMerged)
	partFoot.WriteTo(&partMerged)

	if _, err = w.Write(partMerged.Bytes()); err != nil {
		return err
	}

	return nil
}
