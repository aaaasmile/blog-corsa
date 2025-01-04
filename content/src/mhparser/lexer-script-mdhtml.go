package mhparser

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"text/template"
	"unicode"
)

func lexStateMdHtmlOnLine(l *L) StateFunc {
	for {
		if nextSt, ok := lexMatchFnKey(l); ok {
			return nextSt
		}
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
func lexStateAfterString(l *L) StateFunc {
	for {
		switch r := l.next(); {
		case r == EOFRune || r == '\r' || r == '\n':
			return l.errorf("[lexStateAfterString] Expected next param or close curl")
		case r == ',':
			l.emit(itemSeparator)
			return lexStateInParamString
		case r == ']':
			l.emit(itemEndOfBlock)
			return lexStateMdHtmlBeforeStm
		case r == '\'':
			l.emit(itemParamString)
		case unicode.IsSpace(r):
			l.ignore()
		default:
			return l.errorf("[lexStateAfterString] Malformed end of parameter: %s", l.source[l.start:l.position])
		}
	}
}

func lexStateInParamString(l *L) StateFunc {
	ll := 0
	for {
		rleos := l.peek()
		//fmt.Println("***> ", rleos, ll)
		if ll == 0 && rleos == '\'' {
			l.emit(itemEmptyString)
			return lexStateAfterString
		}

		switch r := l.next(); {
		case r == EOFRune || r == '\r' || r == '\n':
			return l.errorf("[lexStateInParamString] expected end of string")
		case r == '\'':
			l.rewind()
			l.emit(itemParamString)
			return lexStateAfterString
		default:
			ll += 1
		}
	}
}

func lexStateLinkBeforeBegStr(l *L) StateFunc {
	for {
		switch r := l.next(); {
		case unicode.IsSpace(r):
			l.ignore()
		case r == '\r' || r == '\n':
			l.inc_line(r)
			l.ignore()
		case r == '\'':
			l.emit(itemText)
			return lexStateInParamString
		default:
			return l.errorf("[lexStateBeforeCurl] Expected ( but got %s (Line %d)", l.source[l.start:l.position], l.scriptLine)
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

// -- common subtype
type mdhtLineNode struct {
	line string
}

func (n mdhtLineNode) String() string {
	return n.line
}

// -- derived type Link Simple
type mdhtLinkSimpleNode struct {
	mdhtLineNode
	before_link string
	href_arg    string
}

func (ln mdhtLinkSimpleNode) Transform(templDir string) error {
	templName := path.Join(templDir, "transform.html")
	tmplPage := template.Must(template.New("Link").ParseFiles(templName))
	CtxFirst := struct {
		HrefLink    string
		DisplayLink string
	}{
		HrefLink:    ln.href_arg,
		DisplayLink: ln.href_arg,
	}
	var partFirst bytes.Buffer
	if err := tmplPage.ExecuteTemplate(&partFirst, "linkbase", CtxFirst); err != nil {
		return err
	}

	res := fmt.Sprintf("%s%s", ln.before_link, partFirst.String())
	ln.line = res
	return nil
}

// -- Interfaces for Polymorphism
type IMdhtmlLineNode interface {
	String() string
}

type IMdhtmlTransfNode interface {
	IMdhtmlLineNode
	Transform(templDir string) error
}

type MdHtmlGram struct {
	Nodes       []IMdhtmlLineNode
	_curr_Node  IMdhtmlLineNode
	isMdHtmlCtx bool
	debug       bool
}

func NewMdHtmlGr(debug bool) *MdHtmlGram {
	item := MdHtmlGram{
		Nodes: make([]IMdhtmlLineNode, 0),
		debug: debug,
	}
	return &item
}

func (mh *MdHtmlGram) processItem(ll *L, item Token) (bool, error) {
	if item.Type == itemBegMdHtml {
		mh.isMdHtmlCtx = true
		return true, nil
	}
	if !mh.isMdHtmlCtx {
		return false, nil
	}
	switch {
	case item.Type == itemMdHtmlBlockLine:
		mh.Nodes = append(mh.Nodes, mdhtLineNode{line: item.Value})
	case item.Type == itemMdHtmlBlock:
		mh.Nodes = append(mh.Nodes, mdhtLineNode{line: item.Value})
	case item.Type == itemBuiltinFunction:
		if !isLexfnKey(ll, item.ID) {
			return false, fmt.Errorf("[MdHtmlGram] function %s is not defined", item.Value)
		}
		// TODO recognize link node
		mh._curr_Node = mdhtLinkSimpleNode{}
	case item.Type == itemEOF:
		return false, nil
	default:
		return false, fmt.Errorf("[MdHtmlGram] unsupported statement parser %q", item)
	}
	return true, nil
}

func (mh *MdHtmlGram) storeMdHtmlStatement(nrmPrg *NormPrg, scrGr *ScriptGrammar) error {
	if mh.debug {
		fmt.Println("*** storeMdHtmlStatement ", len(mh.Nodes))
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
	for _, node := range mh.Nodes {
		trans, ok := node.(IMdhtmlTransfNode)
		if ok {
			if err := trans.Transform(scrGr.TemplDir); err != nil {
				return err
			}
			linesParam.ArrayValue = append(linesParam.ArrayValue, trans.String())
		} else {
			linesParam.ArrayValue = append(linesParam.ArrayValue, node.String())
		}
	}

	nrmPrg.FnsList = append(nrmPrg.FnsList, fnStMdHtml)
	nrm_st_name, err := nrmPrg.statementInNormMap(stName, scrGr, len(nrmPrg.FnsList)-1)
	if mh.debug {
		fmt.Println("*** storeMdHtmlStatement norm name", nrm_st_name)
	}
	return err
}

func lexMatchFnKey(l *L) (StateFunc, bool) {
	for _, v := range l.descrFns {
		k := fmt.Sprintf("[%s", v.KeyName)
		if strings.HasPrefix(l.source[l.position:], k) { // make sure to parse the longest keyword first
			l.position += len(k)
			l.emitCustFn(v.ItemTokenType, v.CustomID)
			return lexStateLinkBeforeBegStr, true
		}
	}
	return nil, false
}
