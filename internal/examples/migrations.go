package examples

import (
	"gorm.io/gorm"
	"time"
)

func Migrate(db *gorm.DB) error {
	return db.Migrator().AutoMigrate(&ExampleTable{})
}

type ExampleTable struct {
	ID        string     `gorm:"column:id;primaryKey"`
	Desc      string     `gorm:"column:description"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:updated_at"`
}
