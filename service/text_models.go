package service

// --- comment ---

// Params of adding text comment
type AddTextCommentParams struct {
	TextId  int32  `json:"textId" binding:"required,number,gt=0"`
	Comment string `json:"comment" binding:"required,max=200"`
}
