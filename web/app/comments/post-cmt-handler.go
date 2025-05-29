package comments

import (
	"bytes"
	"corsa-blog/conf"
	"corsa-blog/db"
	"corsa-blog/idl"
	"corsa-blog/util"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"text/template"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

type CommentHandler struct {
	debug       bool
	liteCmtDB   *db.LiteDB
	newCmt      chan *idl.CmtItem
	start       time.Time
	moderateCmt bool
}

func NewPostCommentHandler(liteDB *db.LiteDB, debug bool, moderateCmt bool, newCmt chan *idl.CmtItem) *CommentHandler {
	res := &CommentHandler{
		debug:       debug,
		liteCmtDB:   liteDB,
		moderateCmt: moderateCmt,
		start:       time.Now(),
		newCmt:      newCmt,
	}
	return res
}

func (ch *CommentHandler) HandleFormDeleteComment(w http.ResponseWriter, req *http.Request, id int, post_id string) error {
	reqId := req.URL.Query().Get("reqId")
	if reqId == "" {
		return fmt.Errorf("request id for delete is null")
	}
	log.Println("[HandleFormDeleteComment] delete comment ", id, post_id, reqId)
	cmtItem := &idl.CmtItem{
		Id:     id,
		PostId: post_id,
		ReqId:  reqId,
	}

	if err := ch.liteCmtDB.DeleteComment(cmtItem); err != nil {
		return err
	}

	return ch.renderDeletedCmtIdOk(cmtItem, w)
}

func (ch *CommentHandler) HandleFormNewComment(w http.ResponseWriter, req *http.Request, parent_id int, post_id string) error {
	lang := req.URL.Query().Get("lang")
	log.Println("[HandleFormNewComment] process new comment ", parent_id, post_id, lang)
	err := req.ParseForm()
	if err != nil {
		return err
	}
	email := req.PostFormValue("email")
	name := req.PostFormValue("name")
	commentMd := req.PostFormValue("comment")
	dateTimeIso := req.PostFormValue("date")
	if ch.debug {
		log.Println("orig comment:", commentMd)
		log.Println("name, email, date:", name, email, dateTimeIso)
	}
	dtCmt := time.Now()
	if dateTimeIso != "" {
		arr := strings.Split(dateTimeIso, " ")
		parsStr := "2006-01-02"
		if len(arr) == 2 {
			parsStr = "2006-01-02 15:00"
		}
		dt, err := time.Parse(parsStr, dateTimeIso)
		if err != nil {
			return err
		}
		dtCmt = dt
	}
	unsafeComment := blackfriday.Run([]byte(commentMd), blackfriday.WithNoExtensions())
	htmlCmt := bluemonday.StrictPolicy().SanitizeBytes(unsafeComment)
	if ch.debug {
		log.Println("transformed html comment:", string(htmlCmt))
	}

	errMsg := ""
	cmtItem := &idl.CmtItem{
		Email:    email,
		Name:     name,
		Status:   idl.STCreated,
		DateTime: dtCmt,
		Comment:  string(htmlCmt),
		PostId:   post_id,
		ParentId: parent_id,
	}
	if name == "" {
		if _, err := mail.ParseAddress(email); err != nil {
			errMsg = "il nome è vuoto"
			return ch.renderResNewComment(cmtItem, errMsg, w)
		}
	}
	if len(name) > 30 {
		errMsg = "nome non valido"
		return ch.renderResNewComment(cmtItem, errMsg, w)
	}
	if email == "" {
		if conf.Current.Comment.AllowEmptyMail {
			email = "noreply@invido.it"
		} else {
			errMsg = "email è vuota"
			return ch.renderResNewComment(cmtItem, errMsg, w)
		}
	}
	if _, err := mail.ParseAddress(email); err != nil {
		errMsg = "inidirizzo email non valido"
		return ch.renderResNewComment(cmtItem, errMsg, w)
	}
	if len(htmlCmt) == 0 {
		errMsg = "commento vuoto"
		return ch.renderResNewComment(cmtItem, errMsg, w)
	}

	if !ch.moderateCmt {
		cmtItem.Status = idl.STPublished
	}
	cmtItem.ReqId, err = util.PseudoUuid()
	if err != nil {
		return err
	}
	if err := ch.liteCmtDB.InsertNewComment(cmtItem); err != nil {
		return err
	}
	ch.newCmt <- cmtItem

	return ch.renderResNewComment(cmtItem, errMsg, w)
}

func (ch *CommentHandler) renderResNewComment(cmtItem *idl.CmtItem, errMsg string, w http.ResponseWriter) error {
	ctx := struct {
		Cmt       *idl.CmtItem
		ErrMsg    string
		HasErrors bool
		Id        int
		ReqId     string
		ParentId  int
		PostId    string
		PostURL   string
	}{
		Cmt:       cmtItem,
		ErrMsg:    errMsg,
		HasErrors: (errMsg != ""),
		Id:        cmtItem.Id,
		ParentId:  cmtItem.ParentId,
		PostId:    cmtItem.PostId,
		ReqId:     cmtItem.ReqId,
		PostURL:   cmtItem.GetLocationFromPostId(),
	}
	//fmt.Println("*** ctx: ", *ctx.Cmt)

	templName := "templates/cmt/resp-newcomment.html"
	var partMerged bytes.Buffer
	tmplBody := template.Must(template.New("Body").ParseFiles(templName))
	if err := tmplBody.ExecuteTemplate(&partMerged, "base", ctx); err != nil {
		return err
	}

	elapsed := time.Since(ch.start)

	log.Printf("Service total call duration: %v\n", elapsed)
	_, err := w.Write(partMerged.Bytes())
	//fmt.Println("response: ", partMerged.String())
	return err
}

func (ch *CommentHandler) renderDeletedCmtIdOk(cmtItem *idl.CmtItem, w http.ResponseWriter) error {
	ctx := struct {
		Cmt     *idl.CmtItem
		PostURL string
	}{
		Cmt:     cmtItem,
		PostURL: cmtItem.GetLocationFromPostId(),
	}

	templName := "templates/cmt/resp-newcomment.html"
	var partMerged bytes.Buffer
	tmplBody := template.Must(template.New("Body").ParseFiles(templName))
	if err := tmplBody.ExecuteTemplate(&partMerged, "deleteok", ctx); err != nil {
		return err
	}

	elapsed := time.Since(ch.start)

	log.Printf("Service total call duration: %v\n", elapsed)
	_, err := w.Write(partMerged.Bytes())
	return err
}
