package model

type Category struct {
	BaseModel
	Name        string `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Slug        string `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
}
