package response

import "time"

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
    CoverUrl     *string `json:"cover_url"`
    BorrowCount  int     `json:"borrow_count"`
    CreatedAt    string  `json:"created_at"`
}

type FailedItems struct {
    Index int    `json:"index"`
    ISBN  string `json:"isbn"`
    Error string `json:"error"`
}

type BatchCreateBookResponse struct {
    SuccessCount int `json:"success_count"`
    FailedCount  int `json:"failed_count"`
    FailedItems  []FailedItems `json:"failed_items"`
}

type BookListItem struct {
    ID           uint64   `json:"id"`
    Title        string   `json:"title"`
    Author       string   `json:"author"`
    ISBN         string   `json:"isbn"`
    CategoryID   uint     `json:"category_id"`
    CategoryName string   `json:"category_name"`
    Publisher    string   `json:"publisher"`

    PublishDate  *time.Time  `json:"publish_date"`
    Price        *float64 `json:"price"`

    Stock        int      `json:"stock"`
    Available    int      `json:"available"`

    CoverURL     string   `json:"cover_url"`
    BorrowCount  int      `json:"borrow_count"`
}

type GetBookListResponse struct {
    Total      int64          `json:"total"`
    Page       int            `json:"page"`
    Limit      int            `json:"limit"`
    TotalPages int            `json:"total_pages"`
    Books      []BookListItem `json:"books"`
}

type CategoryDetails struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}

type GetBookDetailsResponse struct {
    ID           uint              `json:"id"`
    Title        string            `json:"title"`
    Author       string            `json:"author"`
    ISBN         string            `json:"isbn"`
    Category     CategoryDetails   `json:"category"`
    Publisher    string            `json:"publisher"`
    PublishDate  *time.Time        `json:"publish_date"`
    Price        float64           `json:"price"`
    Stock        int               `json:"stock"`
    Available    int               `json:"available"`
    Description  string            `json:"description"`
    CoverURL     string            `json:"cover_url"`
    BorrowCount  int               `json:"borrow_count"`
    Rating       float64           `json:"rating"`
    CreatedAt    time.Time         `json:"created_at"`
    UpdatedAt    time.Time         `json:"updated_at"`
}

type UpdateBookResponse struct {
    ID           uint64  `json:"id"`
    UpdatedAt    string  `json:"created_at,omitempty"`
    Title       *string  `json:"title,omitempty"`
    Author      *string  `json:"author,omitempty"`
    ISBN        *string  `json:"isbn,omitempty"`
    CategoryID  *uint    `json:"category_id,omitempty"`
    CategoryName *string  `json:"category_name,omitempty"`
    Publisher   *string  `json:"publisher,omitempty"`
    PublishDate *string  `json:"publish_date,omitempty"`
    Price       *float64 `json:"price,omitempty"`
    Stock       *int     `json:"stock,omitempty"`
    Available   *int     `json:"available,omitempty"`
    Description *string  `json:"description,omitempty"`
    CoverURL    *string  `json:"cover_url,omitempty"`
}