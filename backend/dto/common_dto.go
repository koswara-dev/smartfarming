package dto

type Meta struct {
	CurrentPage  int   `json:"currentPage"`
	Limit        int   `json:"limit"`
	TotalRecords int64 `json:"totalRecords"`
	TotalPages   int   `json:"totalPages"`
}

type PaginationResponse struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

type PaginationRequest struct {
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=10"`
	Search string `form:"search"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}
