package idl

import (
	"fmt"
	"path"
	"strings"
	"time"
)

var (
	Appname = "blog-corsa"
	Buildnr = "00.001.20241226-00"
)

type StatusType int

const (
	Created StatusType = iota
	Published
	Deleted
	Rejected
	Spam
)

type CmtItem struct {
	Id       string
	ParentId string
	Name     string
	Email    string
	Website  string
	Comment  string
	Time     time.Time
	Status   StatusType
	Indent   int
}

type CmtNode struct {
	PostId   string
	NodeId   string
	Children []*CmtNode
	CmtItem  *CmtItem
}

type ImgDataItem struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Redux   string `json:"redux"`
	Caption string `json:"caption"`
}

type ImgDataItems struct {
	Images []ImgDataItem `json:"images"`
}

func (fg *ImgDataItem) CalcReduced() error {
	ext := path.Ext(fg.Name)
	if ext == "" {
		return fmt.Errorf("[calcReduced] extension on %s is empty, this is not supported", fg.Name)
	}
	bare_name := strings.Replace(fg.Name, ext, "", -1)
	fg.Redux = fmt.Sprintf("%s_320%s", bare_name, ext)
	return nil
}
