package service

import (
	"strings"

	"github.com/go-playground/validator/v10"
	sttext "github.com/zlbenjamin/gotextgin/pkg/text"
)

// --- comment ---

// Params of adding text
type AddTextParams struct {
	// Text content
	Content string `json:"content" binding:"required,max=10000"`
	// Text type
	Type string `json:"type" binding:"max=10"`
	// Up to 5 tags
	Tags []string `json:"tags" binding:"checktags"`
}

// check tags
// length <= 5, duplicate, blank tag
func CheckTags(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().([]string)
	if ok {
		if len(data) > 5 {
			return false
		}

		tagMap := make(map[string]bool, len(data))
		for _, item := range data {
			item = strings.Trim(item, " ")
			if item == "" {
				return false
			}
			if tagMap[item] {
				return false

			} else {
				tagMap[item] = true
			}
		}
	}

	return true
}

// Convert params to record
func (params AddTextParams) ConvertToText() sttext.Text {
	var ret sttext.Text
	ret.Content = params.Content
	ret.Type = params.Type
	return ret
}

// Get text by ID
type GetTextDTO struct {
	// primary key of text
	Id int32 `json:"id" uri:"id" binding:"required,number,gt=0"`
}

// Params of page find
type pageFindParams struct {
	// page no, start from 1
	PageNo int32 `json:"pageNo" binding:"number,gt=0" example:"1"`
	// page size, range [1, 500]
	PageSize int32 `json:"pageSize" binding:"number,gte=1,lte=500" example:"10"`
	// key word in the content field
	KwContent string `json:"kwContent" binding:"max=50"`
	// type of a text
	Type string `json:"type" binding:"max=10" example:"code"`
	// up to 5 tags
	Tags []string `json:"tags" binding:"checktags" example:"golang,web"`
}

type PageFindVO struct {
	// Text
	sttext.Text
	// Tags of text
	Tags []sttext.TextTag `json:"tags"`
	// Comments of text
	Comments []sttext.TextComment `json:"comments"`
	// Total of comments
	TotalOfComments int64 `json:"totalOfComments" example:"0"`
}

type TextFullVO struct {
	// Text
	sttext.Text
	// Tags of text
	// order by create_time ASC
	Tags []sttext.TextTag `json:"tags"`
	// Comments of text
	// order by create_time ASC
	Comments []sttext.TextComment `json:"comments"`
}

// --- comment ---

// Params of adding text comment
type AddTextCommentParams struct {
	// primary key of a text
	TextId int32 `json:"textId" binding:"required,number,gt=0"`
	// content of the comment
	Comment string `json:"comment" binding:"required,max=200"`
}

// Params of deleting comment
type DelTextCommentParams struct {
	// primary key of text
	TextId int32 `json:"textId" uri:"textId" binding:"required,number,gt=0"`
	// primary key of comment
	Id uint64 `json:"id" uri:"id" binding:"required,number,gt=0"`
}

type GetCommentsDTO struct {
	TextIdList []int32 `json:"textIdList"`
	MaxNumber  int32   `json:"maxNumber"`
}

type GetCommentsVO struct {
	List []TextCommentsVO `json:"list"`
}

type TextCommentsVO struct {
	TextId   int32                `json:"textId"`
	Total    uint64               `json:"total"`
	Comments []sttext.TextComment `json:"comments"`
}
