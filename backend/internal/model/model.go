package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"`
	Nickname  string         `json:"nickname" gorm:"size:50"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:100"`
	Phone     string         `json:"phone" gorm:"uniqueIndex;size:20"`
	Role      string         `json:"role" gorm:"size:20;default:user"`
	Status    int            `json:"status" gorm:"default:1"`
	LastLogin *time.Time     `json:"last_login"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

type Image struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"index;not null"`
	Title       string         `json:"title" gorm:"size:200"`
	Description string         `json:"description" gorm:"type:text"`
	URL         string         `json:"url" gorm:"size:500;not null"`
	Thumbnail   string         `json:"thumbnail" gorm:"size:500"`
	Width       int            `json:"width"`
	Height      int            `json:"height"`
	Size        int64          `json:"size"`
	Format      string         `json:"format" gorm:"size:20"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Image) TableName() string {
	return "images"
}

type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:50;not null"`
	ParentID  uint           `json:"parent_id" gorm:"default:0"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Category) TableName() string {
	return "categories"
}
