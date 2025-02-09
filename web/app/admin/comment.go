package admin

import (
	"corsa-blog/idl"
	"encoding/json"
	"fmt"
)

type CommentParam struct {
	Cmd  string `json:"cmd"`
	Id   string `json:"id"`
	Type string `json:"type"`
}

type CommentReq struct {
	Params CommentParam
}

func (ah *AdminHandler) doComment() error {
	//fmt.Println("*** raw", string(ah.rawbody))
	cmtReq := CommentReq{}
	if err := json.Unmarshal(ah.rawbody, &cmtReq); err != nil {
		return err
	}
	var err error
	cmtPara := cmtReq.Params
	switch cmtPara.Cmd {
	case "list":
		err = ah.doCmtList(cmtPara.Type)
	default:
		return fmt.Errorf("[doComment]comment command not supported %v", cmtPara)
	}
	if err != nil {
		return err
	}
	return nil
}

func (ah *AdminHandler) doCmtList(tt string) error {
	switch tt {
	case "to_moderate":
		if err := ah.doCmtListToModerate(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("[doCmtList] type %s not suported", tt)
	}
	return nil
}

func (ah *AdminHandler) doCmtListToModerate() error {
	cmts, err := ah.liteDB.GeCommentsToModerate()
	if err != nil {
		return err
	}
	resp := struct {
		Comments []*idl.CmtItem
	}{
		Comments: cmts,
	}
	return writeResponse(ah._w, resp)
}
