package mhparser

import (
	"errors"
	"fmt"
	"strings"
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
			return lexStateSingleBeforeBegStr
		case r == ']':
			l.emit(itemEndOfBlock)
			return lexStateMdHtmlOnLine
		case r == '\'':
			l.emit(itemText)
		case unicode.IsSpace(r):
			l.ignore()
		default:
			return l.errorf("[lexStateAfterString] Malformed end of parameter: %s", l.source[l.start:l.position])
		}
	}
}

func lexStateMultiAfterString(l *L) StateFunc {
	for {
		switch r := l.next(); {
		case r == EOFRune || r == '\r' || r == '\n':
			l.ignore()
			l.inc_line(r)
		case r == ',':
			l.emit(itemSeparator)
			return lexStateMultiBeforeBegStr
		case r == ']':
			l.emit(itemEndOfBlock)
			return lexStateMdHtmlOnLine
		case r == '\'':
			l.emit(itemText)
		case unicode.IsSpace(r):
			l.ignore()
		default:
			return l.errorf("[lexStateMultiAfterString] Malformed end of parameter: %s", l.source[l.start:l.position])
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

func lexStateMultiInParamString(l *L) StateFunc {
	ll := 0
	for {
		rleos := l.peek()
		//fmt.Println("***> ", rleos, ll)
		if ll == 0 && rleos == '\'' {
			l.emit(itemEmptyString)
			return lexStateMultiAfterString
		}

		switch r := l.next(); {
		case r == EOFRune || r == '\r' || r == '\n':
			return l.errorf("[lexStateMultiInParamString] expected end of string inside a single line")
		case r == '\'':
			l.rewind()
			l.emit(itemParamString)
			return lexStateMultiAfterString
		default:
			ll += 1
		}
	}
}

func lexStateSingleBeforeBegStr(l *L) StateFunc {
	for {
		switch r := l.next(); {
		case unicode.IsSpace(r):
			l.ignore()
		case r == '\r' || r == '\n':
			return l.errorf("[lexStateSingleBeforeBegStr] Expected next param string")
		case r == '\'':
			l.emit(itemText)
			return lexStateInParamString
		default:
			return l.errorf("[lexStateSingleBeforeBegStr] Expected ( but got %s (Line %d)", l.source[l.start:l.position], l.scriptLine)
		}
	}
}

func lexStateMultiBeforeBegStr(l *L) StateFunc {
	for {
		switch r := l.next(); {
		case unicode.IsSpace(r):
			l.ignore()
		case r == '\r' || r == '\n':
			l.ignore()
			l.inc_line(r)
		case r == '\'':
			l.emit(itemText)
			return lexStateMultiInParamString
		default:
			return l.errorf("[lexStateMultiBeforeBegStr] Expected ( but got %s (Line %d)", l.source[l.start:l.position], l.scriptLine)
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

//  ---  Interfaces

// -- Basic
type IMdhtmlLineNode interface {
	String() string
}

// Node with transformations
type IMdhtmlTransfNode interface {
	IMdhtmlLineNode
	Transform(templDir string) error
	AddParamString(parVal string) error
	AddblockHtml(val string) error
}

// -- Basic, implements IMdhtmlLineNode
type mdhtLineNode struct {
	line        string
	before_link string
	after_link  string
}

func (n *mdhtLineNode) String() string {
	return n.line
}

/////////// ---- Grammar

type MdHtmlGram struct {
	Nodes       []IMdhtmlLineNode
	_curr_Node  IMdhtmlLineNode
	isMdHtmlCtx bool
	debug       bool
	templDir    string
}

func NewMdHtmlGr(templDir string, debug bool) *MdHtmlGram {
	item := MdHtmlGram{
		Nodes:    make([]IMdhtmlLineNode, 0),
		debug:    debug,
		templDir: templDir,
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
		if err := mh.blockHtmlPart(item.Value); err != nil {
			return false, err
		}
	case item.Type == itemMdHtmlBlock:
		if err := mh.blockHtmlPart(item.Value); err != nil {
			return false, err
		}
	case item.Type == itemLinkSimple:
		mh._curr_Node = NewLinkSimpleNode(item.Value)
	case item.Type == itemFigStack:
		mh._curr_Node = NewFigStackNode(item.Value)
	case item.Type == itemText:
		// ignore
	case item.Type == itemSeparator:
		// ignore
	case item.Type == itemParamString:
		if err := mh.addParameterString(item.Value); err != nil {
			return false, err
		}
	case item.Type == itemEndOfBlock:
		if err := mh.endOfBlock(item.Value); err != nil {
			return false, err
		}
	case item.Type == itemEOF:
		return false, nil
	case item.Type == itemError:
		return false, errors.New(item.Value)
	default:
		return false, fmt.Errorf("[MdHtmlGram] unsupported statement parser %q", item)
	}
	return true, nil
}

func (mh *MdHtmlGram) blockHtmlPart(val string) error {
	transf, ok := mh._curr_Node.(IMdhtmlTransfNode)
	if ok {
		if err := transf.AddblockHtml(val); err != nil {
			return err
		}
		mh._curr_Node = &mdhtLineNode{line: "undef"}
	} else {
		mh.Nodes = append(mh.Nodes, &mdhtLineNode{line: val})
	}
	return nil
}

func (mh *MdHtmlGram) addParameterString(valPar string) error {
	trans, ok := mh._curr_Node.(IMdhtmlTransfNode)
	if ok {
		return trans.AddParamString(valPar)
	}
	return fmt.Errorf("param string is not supported")
}

func (mh *MdHtmlGram) endOfBlock(valPar string) error {
	if valPar != "]" {
		return fmt.Errorf("[endOfBlock] end of block not  recognized")
	}
	mh.Nodes = append(mh.Nodes, mh._curr_Node)
	return nil
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
			if v.IsMultiline {
				return lexStateMultiBeforeBegStr, true
			}
			return lexStateSingleBeforeBegStr, true
		}
	}
	return nil, false
}
