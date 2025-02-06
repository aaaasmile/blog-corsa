package admin

import "fmt"

type CommentParam struct {
	Cmd string `json:"cmd"`
	Id  string `json:"id"`
}

type CommentReq struct {
	Params CommentParam
}

func (ah *AdminHandler) doComment() error {
	fmt.Println("*** raw", string(ah.rawbody))

	return nil
}
