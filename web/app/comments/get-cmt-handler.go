package comments

import (
	"bytes"
	"corsa-blog/conf"
	"corsa-blog/db"
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
	if err := tmplBody.ExecuteTemplate(&partForm, "headformDet", ctxHead); err != nil {
		return err
	}

	partForm.WriteTo(&partMerged)

	if _, err = w.Write(partMerged.Bytes()); err != nil {
		return err
	}

	return nil
}

func (ch *CommentHandler) HandleCommentsTitle(w http.ResponseWriter, req *http.Request, post_id string) error {
	lang := req.URL.Query().Get("lang")
	log.Printf("[HandleComments] get comments for id=%s, lang=%s", post_id, lang)

	cmtNode, err := ch.liteDB.GetCommentsForPostId(post_id)
	if err != nil {
		return err
	}

	templName := "templates/cmt/get-comments.html"
	var partHeader bytes.Buffer
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

	if err := tmplBody.ExecuteTemplate(&partHeader, "headTitle", ctxHead); err != nil {
		return err
	}

	if _, err = w.Write(partHeader.Bytes()); err != nil {
		return err
	}

	return nil
}

func (ch *CommentHandler) HandleCommentsDetails(w http.ResponseWriter, req *http.Request, post_id string) error {
	lang := req.URL.Query().Get("lang")
	log.Printf("[HandleCommentsDetails] get comments for id=%s, lang=%s", post_id, lang)

	cmtNode, err := ch.liteDB.GetCommentsForPostId(post_id)
	if err != nil {
		return err
	}

	templName := "templates/cmt/get-comments.html"
	var partForm, partTree, partMerged bytes.Buffer
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

	if err := tmplBody.ExecuteTemplate(&partForm, "headformDet", ctxHead); err != nil {
		return err
	}

	ctxTree := struct {
		CmtLines []string
	}{
		CmtLines: cmtNode.GetLines(),
	}
	if err := tmplBody.ExecuteTemplate(&partTree, "treeDet", ctxTree); err != nil {
		return err
	}

	partTree.WriteTo(&partMerged)
	partForm.WriteTo(&partMerged)

	if _, err = w.Write(partMerged.Bytes()); err != nil {
		return err
	}

	return nil
}
