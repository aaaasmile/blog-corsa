package comments

import (
	"bytes"
	"corsa-blog/conf"
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

func (ch *CommentHandler) HandleFormForReplyComment(w http.ResponseWriter, req *http.Request, id string) error {
	log.Printf("[HandleFormForReplyComment] get Form for comment id=%s", id)

	cmtNode, err := ch.liteDB.GetCommentForId(id)
	if err != nil {
		return err
	}

	templName := "templates/cmt/get-comments.html"
	var partForm, partMerged bytes.Buffer
	tmplBody := template.Must(template.New("DocPart").ParseFiles(templName))

	ctxHead := struct {
		ParentId   int
		PostId     string
		CmtTotText string
		HasDate    bool
	}{
		PostId:   cmtNode.PostId,
		ParentId: cmtNode.CmtItem.Id, // Remember this is a reply form to this Id
		HasDate:  conf.Current.Comment.HasDateInCmtForm,
	}
	if err := tmplBody.ExecuteTemplate(&partForm, "headform", ctxHead); err != nil {
		return err
	}

	partForm.WriteTo(&partMerged)

	if _, err = w.Write(partMerged.Bytes()); err != nil {
		return err
	}

	return nil
}

func (ch *CommentHandler) HandleComments(w http.ResponseWriter, req *http.Request, post_id string) error {
	lang := req.URL.Query().Get("lang")
	log.Printf("[HandleComments] get comments for id=%s, lang=%s", post_id, lang)

	cmtNode, err := ch.liteDB.GeCommentsForPostId(post_id)
	if err != nil {
		return err
	}

	templName := "templates/cmt/get-comments.html"
	var partHeader, partForm, partTree, partFoot, partMerged bytes.Buffer
	tmplBody := template.Must(template.New("DocPart").ParseFiles(templName))

	ctxHead := struct {
		ParentId   int
		PostId     string
		CmtTotText string
		HasDate    bool
	}{
		CmtTotText: cmtNode.GetTextNumComments(),
		PostId:     post_id,
		ParentId:   cmtNode.CmtItem.ParentId,
		HasDate:    conf.Current.Comment.HasDateInCmtForm,
	}
	if err := tmplBody.ExecuteTemplate(&partHeader, "head", ctxHead); err != nil {
		return err
	}
	if err := tmplBody.ExecuteTemplate(&partForm, "headform", ctxHead); err != nil {
		return err
	}

	ctxTree := struct {
		CmtLines []string
	}{
		CmtLines: cmtNode.GetLines(),
	}
	if err := tmplBody.ExecuteTemplate(&partTree, "tree", ctxTree); err != nil {
		return err
	}

	cmtItem := idl.CmtItem{}
	if err := tmplBody.ExecuteTemplate(&partFoot, "foot", cmtItem); err != nil {
		return err
	}
	partHeader.WriteTo(&partMerged)
	partForm.WriteTo(&partMerged)
	partTree.WriteTo(&partMerged)
	partFoot.WriteTo(&partMerged)

	if _, err = w.Write(partMerged.Bytes()); err != nil {
		return err
	}

	return nil
}
