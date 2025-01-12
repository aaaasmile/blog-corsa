package trans

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
type MdhtLineNode struct {
	line        string
	before_link string
	after_link  string
}

func (n *MdhtLineNode) String() string {
	return n.line
}

func NewMdhtLineNode(line string) *MdhtLineNode {
	ln := MdhtLineNode{line: line}
	return &ln
}
