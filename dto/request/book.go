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
