package model

import (
	"time"
)

type Book struct {
	ID			uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title		string    `json:"title" gorm:"type:varchar(200);not null;index:idx_title"`
	Author     	string 	  `json:"author" gorm:"type:varchar(100);not null;index:idx_author"`
	ISBN        string    `json:"isbn" gorm:"type:varchar(20);unique;not null;index:idx_isbn"`
    CategoryID  uint      `json:"category_id" gorm:"index:idx_category;not null"`
    Publisher   string    `json:"publisher" gorm:"type:varchar(100);not null"`
    PublishDate *time.Time
    Price       float64   `json:"price" gorm:"type:decimal(10,2)"`
    Stock       int       `json:"stock" gorm:"default:0"`
    Description string    `json:"description" gorm:"type:text"`
    CoverURL    string    `json:"cover_url" gorm:"type:varchar(500)"`
    BorrowCount int       `json:"borrow_count" gorm:"default:0"`
    Rating      float64   `json:"rating" gorm:"type:decimal(3,2)"`
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

    Category    Category    `gorm:"foreignKey:CategoryID"`
}

type Category struct {
    ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    Name        string    `json:"name" gorm:"type:varchar(50);not null"`
    Description string    `json:"description" gorm:"type:varchar(200)"`
    ParentID    *uint
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

