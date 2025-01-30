package idl

import (
	"fmt"
	"testing"
)

func TestNode(t *testing.T) {
	cI1 := CmtItem{Name: "luz1", Comment: "<p>Hallo 1.0</p>"}
	cn1 := CmtNode{CmtItem: &cI1}
	cI2 := CmtItem{Name: "luz2", Comment: "<p>Hallo 2.0</p>"}
	cn2 := CmtNode{CmtItem: &cI2}

	root := CmtNode{Children: []*CmtNode{&cn1, &cn2}}
	lines := root.GetLines()
	fmt.Println("** lines", lines)
	if len(lines) != 10 {
		t.Error("expected 10 lines, but ", lines)
	}
	if lines[2] != "<p><strong>luz1</strong></p>" {
		t.Error("unexpected line 3: ", lines[2])
	}
}
