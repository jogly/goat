package model

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Base struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	UUID      uuid.UUID  `gorm:"type:uuid" sql:"index" json:"uuid"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
}

func (b *Base) BeforeCreate(s *gorm.Scope) (err error) {
	if b.UUID != uuid.Nil {
		return
	}
	b.UUID, err = uuid.NewV4()
	return
}
