package response

type CreateBookResponse struct {
    ID           uint64  `json:"id"`
    Title        string  `json:"title"`
    Author       string  `json:"author"`
    ISBN         string  `json:"isbn"`
    CategoryID   uint    `json:"category_id"`
    CategoryName string  `json:"category_name"`
    Publisher    string  `json:"publisher"`
    PublishDate  *string `json:"publish_date"`
    Price        *float64 `json:"price"`
    Stock        int     `json:"stock"`
    Available    int     `json:"available"`
    BorrowCount  int     `json:"borrow_count"`
    CreatedAt    string  `json:"created_at"`
}