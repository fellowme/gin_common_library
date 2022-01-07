package param

type PageRequestParam struct {
	Page     int `json:"page,omitempty" form:"page" binding:"required"`
	PageSize int `json:"page_size,omitempty" form:"page_size" binding:"required"`
}
