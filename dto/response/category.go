package response

import "time"

type CategoryItem struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	BookCount   int64     `json:"book_count,omitempty"`
	ParentID    *uint     `json:"parent_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetCategoryListResponse struct {
	Categories []CategoryItem `json:"categories"`
}

type CreateCategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	ParentID    *uint     `json:"parent_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type CategoryChild struct {
	ID        uint  `json:"id"`
	Name      string `json:"name"`
	BookCount int64  `json:"book_count"`
}

type GetCategoryDetailResponse struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	BookCount   int64           `json:"book_count"`
	ParentID    *uint           `json:"parent_id"`
	Children    []CategoryChild `json:"children,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

type UpdateCategoryResponse struct {
	ID        uint      `json:"id"`
	Name      *string   `json:"name,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}