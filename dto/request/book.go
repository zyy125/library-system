package request

type CreateBookRequest struct {
    Title       string  `json:"title" binding:"required,min=1,max=200"`
    Author      string  `json:"author" binding:"required,min=1,max=100"`
    ISBN        string  `json:"isbn" binding:"required"`
    CategoryID  uint    `json:"category_id" binding:"required"`
    Publisher   string  `json:"publisher" binding:"required,min=1,max=100"`
    PublishDate *string `json:"publish_date"`
    Price       *float64 `json:"price" binding:"omitempty,gt=0"`
    Stock       int     `json:"stock" binding:"gte=0"`
    Description *string `json:"description" binding:"max=1000"`
    CoverURL    *string `json:"cover_url"`
}

type BatchCreateBookRequest struct {
    Books  []CreateBookRequest `json:"books"`
}

type GetBookListRequest struct {
    Page          int     `form:"page" binding:"omitempty,min=1"`
    Limit         int     `form:"limit" binding:"omitempty,min=1,max=100"`

    Title         *string `form:"title"`
    Author        *string `form:"author"`
    ISBN          *string `form:"isbn"`
    CategoryID    *uint   `form:"category_id"`
    Publisher     *string `form:"publisher"`

    AvailableOnly *bool   `form:"available_only"`

    SortBy        *string `form:"sort_by" binding:"omitempty,oneof=title author publish_date borrow_count created_at"`
    Order         *string `form:"order" binding:"omitempty,oneof=asc desc"`
}

type UpdateBookRequest struct {
    Title       *string  `json:"title" binding:"omitempty,min=1,max=200"`
    Author      *string  `json:"author" binding:"omitempty,min=1,max=100"`
    ISBN        *string  `json:"isbn"`
    CategoryID  *uint    `json:"category_id"`
    Publisher   *string  `json:"publisher" binding:"omitempty,min=1,max=100"`
    PublishDate *string  `json:"publish_date"`
    Price       *float64 `json:"price" binding:"omitempty,gt=0"`
    Stock       *int     `json:"stock" binding:"omitempty,gte=0"`
    Description *string  `json:"description" binding:"omitempty,max=1000"`
    CoverURL    *string  `json:"cover_url"`
}


