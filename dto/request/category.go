package request

type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=50"`
	Description *string `json:"description" binding:"omitempty,max=200"`
	ParentID    *uint   `json:"parent_id"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=50"`
	Description *string `json:"description" binding:"omitempty,max=200"`
	ParentID    *uint   `json:"parent_id"`
}

type GetCategoryListRequest struct {
	IncludeCount bool `form:"include_count"`
}