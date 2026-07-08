package model

type User struct {
	BaseModel
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	Email    string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Role     string `gorm:"type:varchar(50);default:'user';not null" json:"role"`
	PhotoURL string `gorm:"type:varchar(512)" json:"photoUrl"`
}
