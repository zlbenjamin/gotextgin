package service

import (
	"strings"

	"github.com/go-playground/validator/v10"
	sttext "github.com/zlbenjamin/gotextgin/pkg/text"
)

// --- comment ---

// Params of adding text
type AddTextParams struct {
	Content string   `json:"content" binding:"required,max=10000"`
	Type    string   `json:"type" binding:"max=10"`
	Tags    []string `json:"tags" binding:"checktags"`
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

// Params of page find
type pageFindParams struct {
	PageNo    int32    `json:"pageNo" binding:"number,gt=0"`
	PageSize  int32    `json:"pageSize" binding:"number,gte=1,lte=500"`
	KwContent string   `json:"kwContent" binding:"max=50"`
	Type      string   `json:"type" binding:"max=10"`
	Tags      []string `json:"tags" binding:"checktags"`
}

type PageFindVO struct {
	sttext.Text
	Tags []sttext.TextTag `json:"tags"`
}

// --- comment ---

// Params of adding text comment
type AddTextCommentParams struct {
	TextId  int32  `json:"textId" binding:"required,number,gt=0"`
	Comment string `json:"comment" binding:"required,max=200"`
}

// Params of deleting comment
type DelTextCommentParams struct {
	// Error:Field validation for 'Id' failed on the 'required' tag
	// Id uint64 `json:"id" binding:"required,number,gt=0"` // no
	// Error:Field validation for 'Id' failed on the 'gt' tag
	// Id uint64 `json:"id" binding:"number,gt=0"` // no
	Id uint64 `json:"id" uri:"id" binding:"required,number,gt=0"`
}
