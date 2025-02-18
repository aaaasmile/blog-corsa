package admin

import (
	"corsa-blog/idl"
	"encoding/json"
	"fmt"
	"log"
)

type CommentParam struct {
	Cmd  string `json:"cmd"`
	Id   string `json:"id"`
	Type string `json:"type"`
}

type CommentParamList struct {
	Cmd string `json:"cmd"`
	Ids []int  `json:"list"`
}

type CommentReq struct {
	Params CommentParam
}

type CommentReqWithList struct {
	Params CommentParamList
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
	case "approve":
		err = ah.doCmtApprove()
	case "reject":
		err = ah.doCmtReject()
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

func (ah *AdminHandler) doCmtApprove() error {
	cmtReqList := CommentReqWithList{}
	if err := json.Unmarshal(ah.rawbody, &cmtReqList); err != nil {
		return err
	}
	log.Println("[doCmtApprove] the list ", cmtReqList.Params.Ids)
	if err := ah.liteDB.ApproveComments(cmtReqList.Params.Ids); err != nil {
		return err
	}
	//fmt.Println("*** ", string(ah.rawbody))
	return ah.doCmtListToModerate()
}

func (ah *AdminHandler) doCmtReject() error {
	cmtReqList := CommentReqWithList{}
	if err := json.Unmarshal(ah.rawbody, &cmtReqList); err != nil {
		return err
	}
	log.Println("[doCmtReject] the list ", cmtReqList.Params.Ids)
	if err := ah.liteDB.RejectComments(cmtReqList.Params.Ids); err != nil {
		return err
	}
	return ah.doCmtListToModerate()
}

func (ah *AdminHandler) doCmtListToModerate() error {
	cmts, err := ah.liteDB.GetCommentsToModerate()
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
