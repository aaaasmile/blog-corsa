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
	STCreated StatusType = iota
	STPublished
	STDeleted
	STRejected
	STSpam
)

type CmtItem struct {
	Id       int
	ParentId int
	Name     string
	Email    string
	Comment  string
	DateTime time.Time
	Status   StatusType
	Indent   int
	PostId   string
	ReqId    string
}

type CmtNode struct {
	PostId         string
	Children       []*CmtNode
	CmtItem        *CmtItem
	NodeCount      int
	PublishedCount int
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
	fg.Id = fmt.Sprintf("%s_%s", bare_name, fg.Id)
	return nil
}

type PostItem struct {
	Id             int
	PostId         string
	Title          string
	TitleImgUri    string
	DateTime       time.Time
	DateTimeRfC822 string
	Abstract       string
	Uri            string
	Md5            string
}

type PostLinks struct {
	PrevLink   string
	PrevPostID string
	NextLink   string
	NextPostID string
	Item       *PostItem
}

type MapPostsLinks struct {
	MapPost  map[string]PostLinks
	ListPost []PostItem
}
