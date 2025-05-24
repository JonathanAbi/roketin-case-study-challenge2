package entity

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	ID          int             `gorm:"primaryKey" json:"id"`
	Title       string          `gorm:"type:varchar(255); not null" json:"title"`
	Description string          `gorm:"type:text" json:"description"`
	Duration    int             `json:"duration_minutes"`
	Artists     string          `gorm:"type:varchar(255)" json:"artists"`
	Genres      string          `gorm:"type:varchar(255)" json:"genres"`
	FilePath    string          `gorm:"type:varchar(255)" json:"file_path"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty"` //soft delete
}

type MovieFilter struct {
	Title       string
	Description string
	Genres      []string
	Artists     []string
	Page        int
	Limit       int
}

func (Movie) TableName() string {
	return "movies"
}

func (f *MovieFilter) GetPage() int {
	if f.Page <= 0 {
		return 1
	}
	return f.Page
}

func (f *MovieFilter) GetLimit() int {
	if f.Limit <= 0 {
		return 2
	}
	return f.Limit
}
