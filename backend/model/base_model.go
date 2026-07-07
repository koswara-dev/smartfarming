package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	CreatedBy *uuid.UUID     `gorm:"type:uuid" json:"createdBy,omitempty"`
	UpdatedBy *uuid.UUID     `gorm:"type:uuid" json:"updatedBy,omitempty"`
	DeletedBy *uuid.UUID     `gorm:"type:uuid" json:"deletedBy,omitempty"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	if userIDVal := tx.Statement.Context.Value("userID"); userIDVal != nil {
		if userID, ok := userIDVal.(uuid.UUID); ok {
			b.CreatedBy = &userID
			b.UpdatedBy = &userID
		}
	}
	return nil
}

func (b *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	if userIDVal := tx.Statement.Context.Value("userID"); userIDVal != nil {
		if userID, ok := userIDVal.(uuid.UUID); ok {
			b.UpdatedBy = &userID
		}
	}
	return nil
}

func (b *BaseModel) BeforeDelete(tx *gorm.DB) error {
	if userIDVal := tx.Statement.Context.Value("userID"); userIDVal != nil {
		if userID, ok := userIDVal.(uuid.UUID); ok {
			b.DeletedBy = &userID
			tx.Statement.SetColumn("deleted_by", userID)
		}
	}
	return nil
}
