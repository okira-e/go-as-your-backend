package models

import (
	"time"

	"github.com/samborkent/uuidv7"
)

type Post struct {
	ID        string     `sql:"id"             gorm:"type:uuid;primaryKey"`
	Title     string     `sql:"title"          gorm:"size:255;not null"`
	Content   string     `sql:"content"        gorm:"type:text"`
	Published bool       `sql:"published"      gorm:"not null;default:false"`
	CreatedAt time.Time  `sql:"created_at"     gorm:"not null;default:now()"`
	UpdatedAt *time.Time `sql:"updated_at"`
	UserID    string     `sql:"user_id"        gorm:"type:uuid;not null"`

	User User
}

func (self *Post) ToDto() *PostDto {
	return &PostDto{
		ID:        self.ID,
		UserID:    self.UserID,
		Title:     self.Title,
		Content:   self.Content,
		Published: self.Published,
		CreatedAt: self.CreatedAt,
		UpdatedAt: self.UpdatedAt,
	}
}

type CreatePostDto struct {
	Title     string `json:"title"     validate:"required,min=1,max=255"`
	Content   string `json:"content"`
	Published bool   `json:"published"`
}

func (self *CreatePostDto) FromDto(userId string) *Post {
	id := uuidv7.New().String()

	return &Post{
		ID:        id,
		UserID:    userId,
		Title:     self.Title,
		Content:   self.Content,
		Published: self.Published,
	}
}

type PostDto struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"     validate:"required,min=1,max=255"`
	Content   string     `json:"content"`
	Published bool       `json:"published"`
	UserID    string     `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
