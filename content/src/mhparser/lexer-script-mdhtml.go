package mhparser

import (
	"fmt"
	"strings"
)

func lexStateMdHtmlOnLine(l *L) StateFunc {
	for {
		switch r := l.next(); {
		case r == EOFRune:
			l.emit(itemMdHtmlBlock)
			return nil
		case r == '\r' || r == '\n':
			l.rewind()
			l.emit(itemMdHtmlBlockLine)
			l.inc_line(r)
			l.next()
			l.ignore()
		}
	}
}

func lexStateMdHtmlBeforeStm(l *L) StateFunc {
	for {
		switch r := l.next(); {
		case r == EOFRune:
			return nil
		case r == '\r' || r == '\n':
			l.rewind()
			l.emit(itemBegMdHtml)
			l.inc_line(r)
			l.next()
			l.ignore()
			return lexStateMdHtmlOnLine
		case r == '-':
			// nothing
		default:
			return l.errorf("[lexStateMdHtmlBeforeStm] Unexpected char in data separator %s ", l.source[l.start:l.position])
		}
	}
}

func lexMatchMdHtmlKey(l *L) (StateFunc, bool) {
	khtml := "---"
	if strings.HasPrefix(l.source[l.position:], khtml) {
		return lexStateMdHtmlBeforeStm, true
	}
	return nil, false
}

type MdHtmlGram struct {
	Lines       []string
	isMdHtmlCtx bool
	debug       bool
}

func NewMdHtmlGr(debug bool) *MdHtmlGram {
	item := MdHtmlGram{
		Lines: make([]string, 0),
		debug: debug,
	}
	return &item
}

func (mh *MdHtmlGram) processItem(item Token) (bool, error) {
	if item.Type == itemBegMdHtml {
		mh.isMdHtmlCtx = true
		return true, nil
	}
	if !mh.isMdHtmlCtx {
		return false, nil
	}
	switch {
	case item.Type == itemMdHtmlBlockLine:
		mh.Lines = append(mh.Lines, item.Value)
	case item.Type == itemMdHtmlBlock:
		mh.Lines = append(mh.Lines, item.Value)
	case item.Type == itemEOF:
		return false, nil
	default:
		return false, fmt.Errorf("[MdHtmlGram] unsupported statement parser %q", item)
	}
	return true, nil
}

func (mh *MdHtmlGram) storeMdHtmlStatement(nrmPrg *NormPrg, scrGr *ScriptGrammar) error {
	if mh.debug {
		fmt.Println("*** storeMdHtmlStatement ", len(mh.Lines))
	}

	stName := "mdhtml"
	fnStMdHtml := FnStatement{
		IsInternal: true,
		FnName:     stName,
		Type:       TtHtmlVerbatim,
		Params:     make([]ParamItem, 1),
	}
	linesParam := &fnStMdHtml.Params[0]
	linesParam.Label = "Lines"
	linesParam.IsArray = true
	linesParam.ArrayValue = make([]string, 0)
	linesParam.ArrayValue = append(linesParam.ArrayValue, mh.Lines...)

	nrmPrg.FnsList = append(nrmPrg.FnsList, fnStMdHtml)
	nrm_st_name, err := nrmPrg.statementInNormMap(stName, scrGr, len(nrmPrg.FnsList)-1)
	if mh.debug {
		fmt.Println("*** storeMdHtmlStatement norm name", nrm_st_name)
	}
	return err

}
