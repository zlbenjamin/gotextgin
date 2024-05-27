package pkg

import "time"

// table name: text
const Table_Text = "text"

// table name: text_comment
const Table_Text_Comment = "text_comment"

// table name: text_tag
const Table_Text_Tag = "text_tag"

// table: text
type Text struct {
	Id         int32     `json:"id"`
	Content    string    `json:"content"`
	Type       string    `json:"type"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

// table: text_comment
type TextComment struct {
	// TODO
}

// table: text_tag
type TextTag struct {
	// TODO
}
