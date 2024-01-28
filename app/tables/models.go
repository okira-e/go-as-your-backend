package tables

import (
	"time"
)

// ADDED TABLES SHOULD BE ADDED TO app/datasource/datasource.go

type Organizations struct {
	ID          uint      `gorm:"primaryKey; autoIncrement" json:"id"`
	Name        string    `gorm:"unique; not null" validate:"required" json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type Products struct {
	ID          uint      `gorm:"primaryKey; autoIncrement" json:"id"`
	Name        string    `gorm:"unique; not null" validate:"required" json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type Projects struct {
	Product *Products

	ID          uint   `gorm:"primaryKey; autoIncrement" json:"id"`
	Name        string `gorm:"unique; not null" validate:"required" json:"name"`
	Description string `json:"description"`
	// Foreign key
	ProductID uint           `json:"productId"`
	Users     []*SystemUsers `gorm:"many2many:users_projects" json:"users"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
}

type SystemUsers struct {
	ID        uint        `gorm:"primaryKey; autoIncrement" json:"id"`
	FirstName string      `gorm:"not null" validate:"required" json:"firstName"`
	LastName  string      `gorm:"not null" validate:"required" json:"lastName"`
	Email     string      `gorm:"unique; not null" validate:"required,email" json:"email"`
	Title     string      `gorm:"not null" validate:"required" json:"title"`
	Password  string      `gorm:"not null" validate:"required" json:"password"`
	Projects  []*Projects `gorm:"many2many:users_projects; constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"projects"`
	Roles     []*Roles    `gorm:"many2many:users_roles; constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"roles"`
	CreatedAt time.Time   `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime" json:"updatedAt"`
}

type Roles struct {
	ID          uint           `gorm:"primaryKey; autoIncrement" json:"id"`
	Name        string         `gorm:"unique; not null" validate:"required" json:"name"`
	Description string         `gorm:"not null" validate:"required" json:"description"`
	Users       []*SystemUsers `gorm:"many2many:users_roles; constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"users"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
}

