package admin

import (
	"corsa-blog/conf"
	"corsa-blog/crypto"
	"corsa-blog/db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type AdminHandler struct {
	_w      http.ResponseWriter
	_req    *http.Request
	rawbody []byte
	liteDB  *db.LiteDB
}

func NewAdmin(w http.ResponseWriter, req *http.Request, liteDB *db.LiteDB) *AdminHandler {
	ah := AdminHandler{_w: w, _req: req, liteDB: liteDB}
	return &ah
}

func (ah *AdminHandler) HandleAdminRequest() error {
	var err error
	ah.rawbody, err = io.ReadAll(ah._req.Body)
	if err != nil {
		return err
	}

	scopeDef := struct {
		Method string `json:"method"`
	}{}
	if err := json.Unmarshal(ah.rawbody, &scopeDef); err != nil {
		return err
	}
	if scopeDef.Method != "DoLogin" {
		if err := ah.checkReqAuthorization(); err != nil {
			return err
		}
	}
	switch scopeDef.Method {
	case "DoLogin":
		err = ah.doLogin()
	case "DoComment":
		err = ah.doComment()
	default:
		return fmt.Errorf("[HandleAdminRequest]%s is  not supported", scopeDef.Method)
	}
	if err != nil {
		return err
	}
	return nil
}

func (ah *AdminHandler) checkReqAuthorization() error {
	refCred := conf.Current.AdminCred
	req_auth := ah._req.Header.Get("Authorization")
	if req_auth == "" {
		return fmt.Errorf("no authorization provided")
	}
	user, err := refCred.ParseJwtToken(req_auth)
	if err != nil {
		return err
	}
	log.Println("token request from ", user)
	return nil
}

func tokenResult(resultCode int, username string, w http.ResponseWriter) error {
	resp := struct {
		Info       string
		ResultCode int
		Username   string
		Token      crypto.Token
	}{
		ResultCode: resultCode,
		Username:   username,
	}

	switch resultCode {
	case 200:
		resp.Info = "User credential OK"
		expires := 3600 * 24
		log.Printf("Create JWT Token for user %s, expires in %d", username, expires)
		refCred := conf.Current.AdminCred
		err := refCred.GetJWTToken(username, expires, &resp.Token)
		if err != nil {
			return err
		}
		return writeResponse(w, &resp)
	case 403:
		resp.Info = "User Unauthorized"
	default:
		resp.Info = "User credential ERROR"
	}

	return writeErrorResponse(w, resp.ResultCode, resp)
}

func writeResponse(w http.ResponseWriter, resp interface{}) error {
	blobresp, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(blobresp)
	return nil
}

func writeErrorResponse(w http.ResponseWriter, errorcode int, resp interface{}) error {
	blobresp, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	http.Error(w, string(blobresp), errorcode)
	return nil
}
