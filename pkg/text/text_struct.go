package text

import "time"

// table name: text
const Table_Text = "text"

// table name: text_comment
const Table_Text_Comment = "text_comment"

// table name: text_tag
const Table_Text_Tag = "text_tag"

// table: text
type Text struct {
	// Primary key
	Id int32 `json:"id" gorm:"primaryKey;comment:ID"`
	// Text content
	Content string `json:"content" gorm:"not null;comment:content of a Text;type:MEDIUMTEXT"`
	// Text type
	Type string `json:"type" gorm:"not null;size:100;comment:type of a Text"`
	// Create time
	CreateTime time.Time `json:"createTime" gorm:"not null;comment:create time;default:CURRENT_TIMESTAMP;index"`
	// The latest update time
	UpdateTime time.Time `json:"updateTime" gorm:"not null;comment:the latest update time;default:CURRENT_TIMESTAMP"`
}

// table: text_comment
type TextComment struct {
	// Primary key
	Id uint64 `json:"id" gorm:"primaryKey;comment:ID"`
	// Primary key of text
	TextId int32 `json:"textId" gorm:"not null;comment:ID of a Text;index"`
	// Content
	Comment string `json:"comment" gorm:"not null;size:500;comment:content of comment"`
	// Create time
	CreateTime time.Time `json:"createTime" gorm:"not null;comment:create time;default:CURRENT_TIMESTAMP;index"`
	// The latest update time
	UpdateTime time.Time `json:"updateTime" gorm:"not null;comment:the latest update time;default:CURRENT_TIMESTAMP"`
}

// table: text_tag
type TextTag struct {
	// Primary key
	Id uint64 `json:"id" gorm:"primaryKey;comment:ID"`
	// Primary key of text
	TextId int32 `json:"textId" gorm:"not null;comment:ID of a Text;index"`
	// Tag name
	Name string `json:"name" gorm:"not null;size:500;comment:Tag name"`
	// Create time
	CreateTime time.Time `json:"createTime" gorm:"not null;comment:create time;default:CURRENT_TIMESTAMP"`
	// The latest update time
	UpdateTime time.Time `json:"updateTime" gorm:"not null;comment:the latest update time;default:CURRENT_TIMESTAMP"`
}
