package pkg

// Response data
type ApiResponse struct {
	Code    int32  `json:"code" example:"200"`
	Message string `json:"message" example:"OK"`
	Data    any    `json:"data"`
}

// Response data of page find data
// As ApiResponse.Data
type ApiPageFindResponse struct {
	PageNo    int32 `json:"pageNo" example:"1"`
	PageSize  int32 `json:"pageSize" example:"10"`
	Total     int64 `json:"total" example:"0"`
	TotalPage int64 `json:"totalPage" example:"0"`
	List      any   `json:"list"`
}
