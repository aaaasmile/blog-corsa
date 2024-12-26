package idl

import (
	"time"
)

var (
	Appname = "comment-blog"
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
	Id       string
	Children []*CmtNode
	Item     *CmtItem
}
