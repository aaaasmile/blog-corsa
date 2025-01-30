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

func (cn *CmtNode) GetLines() []string {
	res := []string{}
	res = append(res, "<ul>")
	for _, item := range cn.Children {
		if item.CmtItem != nil {
			lines := item.getNodeLines()
			res = append(res, lines...)
		}
	}
	res = append(res, "</ul>")
	return res
}

func (cn *CmtNode) getNodeLines() []string {
	l1 := fmt.Sprintf("<p><strong>%s</strong></p>", cn.CmtItem.Name)
	l2 := fmt.Sprintf("%s<button>Rispondi</button>", cn.CmtItem.Comment)
	res := []string{"<li>", l1, l2}
	if len(cn.Children) > 0 {
		res = append(res, "<ul>")
		for _, item := range cn.Children {
			lines := item.getNodeLines()
			res = append(res, lines...)
		}
		res = append(res, "</ul>")
	}
	res = append(res, "</li>")
	return res
}
