package tables

import (
	"time"
)

// ADDED TABLES SHOULD BE ADDED TO app/datasource/datasource.go

type Organizations struct {
	ID          uint      `gorm:"primaryKey; autoIncrement" json:"id"`
	Name        string    `gorm:"unique; not null" json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type Products struct {
	ID          uint      `gorm:"primaryKey; autoIncrement" json:"id"`
	Name        string    `gorm:"unique; not null" json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type Projects struct {
	Product *Products

	ID          uint   `gorm:"primaryKey; autoIncrement" json:"id"`
	Name        string `gorm:"unique; not null" json:"name"`
	Description string `json:"description"`
	// Foreign key
	ProductID uint      `json:"productId"`
	Users     []*Users  `gorm:"many2many:users_projects" json:"users"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type Users struct {
	ID        uint        `gorm:"primaryKey; autoIncrement" json:"id"`
	FirstName string      `gorm:"not null" json:"firstName"`
	LastName  string      `gorm:"not null" json:"lastName"`
	Email     string      `gorm:"not null" json:"email"`
	Title     string      `gorm:"not null" json:"title"`
	Password  string      `gorm:"not null" json:"password"`
	Projects  []*Projects `gorm:"many2many:users_projects"`
	Roles     []*Roles    `gorm:"many2many:users_roles"`
	CreatedAt time.Time   `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime" json:"updatedAt"`
}

type Roles struct {
	ID          uint      `gorm:"primaryKey; autoIncrement" json:"id"`
	Name        string    `gorm:"unique; not null" json:"name"`
	Description string    `gorm:"not null" json:"description"`
	Users       []*Users  `gorm:"many2many:users_roles"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
