package pkg

// Response data
type ApiResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// Response data of page find data
// As ApiResponse.Data
type ApiPageFindResponse struct {
	PageNo    int32 `json:"pageNo"`
	PageSize  int32 `json:"pageSize"`
	Total     int64 `json:"total"`
	TotalPage int64 `json:"totalPage"`
	List      any   `json:"list"`
}
