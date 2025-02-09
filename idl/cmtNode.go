package idl

import (
	"corsa-blog/util"
	"fmt"
)

func (cn *CmtNode) GetTextNumComments() string {
	if cn.PublishedCount == 1 {
		return "1 commento"
	} else if cn.PublishedCount == 0 {
		return "Nessun commento"
	}
	return fmt.Sprintf("%d commenti", cn.PublishedCount)
}

func (cn *CmtNode) GetLines() []string {
	res := []string{}
	res = append(res, "<ul>")
	for _, item := range cn.Children {
		if item.CmtItem != nil {
			if item.CmtItem.Status == STPublished {
				lines := item.getNodeLines()
				res = append(res, lines...)
			}
		}
	}
	res = append(res, "</ul>")
	return res
}

func (cn *CmtNode) getNodeLines() []string {
	l1 := fmt.Sprintf("<p><strong>%s</strong>, <em><small>%s</small></em></p>", cn.CmtItem.Name, util.FormatDateIt(cn.CmtItem.DateTime))
	l2 := fmt.Sprintf("%s<button hx-get=\"/blog-admin/%d/cmtform\"  hx-target=\"#reply%d\">Rispondi</button><span id=\"reply%d\"></span>", cn.CmtItem.Comment, cn.CmtItem.Id, cn.CmtItem.Id, cn.CmtItem.Id)
	res := []string{"<li>", l1, l2}
	if len(cn.Children) > 0 {
		res = append(res, "<ul>")
		for _, item := range cn.Children {
			if item.CmtItem.Status == STPublished {
				lines := item.getNodeLines()
				res = append(res, lines...)
			}
		}
		res = append(res, "</ul>")
	}
	res = append(res, "</li>")
	return res
}
