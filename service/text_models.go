package service

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
