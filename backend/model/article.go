package model

import (
	"github.com/google/uuid"
)

type Article struct {
	BaseModel
	Title      string    `gorm:"type:varchar(255);not null" json:"title"`
	Slug       string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null" json:"categoryId"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category"`
	ImageURL   string    `gorm:"type:varchar(512)" json:"imageUrl"`
}
