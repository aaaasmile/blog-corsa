package mhparser

import "fmt"

func (i TokenType) String() string {
	switch i {
	case itemText:
		return "itemText"
	case itemBuiltinFunction:
		return "itemBuiltinFunction"
	case itemStringValue:
		return "itemStringValue"
	case itemAssign:
		return "itemAssign"
	case itemComment:
		return "itemComment"
	case itemEmptyString:
		return "itemEmptyString"
	case itemEndOfStatement:
		return "itemEndOfStatement"
	case itemError:
		return "itemError"
	case itemFunctionName:
		return "itemFunctionName"
	case itemFunctionStartBlock:
		return "itemFunctionStartBlock"
	case itemFunctionEnd:
		return "itemFunctionEnd"
	case itemEOF:
		return "itemEOF"
	case itemArrayBegin:
		return "itemArrayBegin"
	case itemArrayEnd:
		return "itemArrayEnd"
	case itemVariable:
		return "itemVariable"
	case itemParamString:
		return "itemParamString"
	}
	return fmt.Sprintf("TokenType %d undef", i)
}
