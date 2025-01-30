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

func TestNodeWithChild(t *testing.T) {
	cI1 := CmtItem{Name: "luz1", Comment: "<p>Hallo 1.0</p>"}
	cn1 := CmtNode{CmtItem: &cI1}
	cI2 := CmtItem{Name: "luz2", Comment: "<p>Hallo 2.0</p>"}
	cn2 := CmtNode{CmtItem: &cI2}

	cI1_1 := CmtItem{Name: "luz11", Comment: "<p>Hallo 1.1</p>"}
	cn1_2 := CmtNode{CmtItem: &cI1_1}
	cn1.Children = append(cn1.Children, &cn1_2)

	root := CmtNode{Children: []*CmtNode{&cn1, &cn2}}
	lines := root.GetLines()
	fmt.Println("** lines", lines)
	if len(lines) != 16 {
		t.Error("expected 10 lines, but ", lines)
	}
	if lines[2] != "<p><strong>luz1</strong></p>" {
		t.Error("unexpected line 3: ", lines[2])
	}
}
