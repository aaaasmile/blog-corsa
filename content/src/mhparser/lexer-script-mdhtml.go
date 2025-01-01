package mhparser

import "strings"

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
