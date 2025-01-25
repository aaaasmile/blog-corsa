package idl

import "fmt"

func (cn *CmtNode) GetTextNumComments() string {
	if cn.NodeCount == 1 {
		return "1 commento"
	} else if cn.NodeCount == 0 {
		return "Nessun commento"
	}
	return fmt.Sprintf("%d commenti", cn.NodeCount)
}
