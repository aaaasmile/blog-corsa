package idl

import "fmt"

func (cn *CmtNode) GetTextNumComments() string {
	if cn.NodeCount == 1 {
		return "1 commento"
	}
	return fmt.Sprintf("%d commenti", cn.NodeCount)
}
