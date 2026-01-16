package models

import (
	"time"

	"github.com/samborkent/uuidv7"
)

type User struct {
	ID string `sql:"id"             gorm:"type:uuid;primaryKey"`
	// RoleID can be null meaning a normal user
	RoleID    *string    `sql:"role_id"        gorm:"type:uuid"`
	FirstName string     `sql:"first_name"     gorm:"size:32;not null"`
	LastName  string     `sql:"last_name"      gorm:"size:32;not null"`
	Email     string     `sql:"email"          gorm:"type:text;uniqueIndex;not null"`
	Password  string     `sql:"password"       gorm:"type:text;not null"`
	Phone     string     `sql:"phone"          gorm:"type:text;not null;unique"`
	IsActive  bool       `sql:"is_active"      gorm:"not null"`
	CreatedAt time.Time  `sql:"created_at"     gorm:"not null;default:now()"`
	UpdatedAt *time.Time `sql:"updated_at"`

	Role  Role
	Posts []Post
}

func (self *User) ToDto() *UserDto {
	return &UserDto{
		ID:        self.ID,
		RoleID:    self.RoleID,
		FirstName: self.FirstName,
		LastName:  self.LastName,
		Email:     self.Email,
		Phone:     self.Phone,
		CreatedAt: self.CreatedAt,
		UpdatedAt: self.UpdatedAt,
	}
}

type UserDto struct {
	ID string `json:"id"`
	// RoleID can be null meaning a normal user
	RoleID    *string    `json:"role_id"`
	FirstName string     `json:"first_name"   validate:"required"`
	LastName  string     `json:"last_name"    validate:"required"`
	Email     string     `json:"email"        validate:"required,email"`
	Phone     string     `json:"phone"        validate:"required,e164"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (self *UserDto) FromDto(autoFill bool, password string) *User {
	id := self.ID
	if autoFill {
		id = uuidv7.New().String()
	}

	return &User{
		ID:        id,
		RoleID:    self.RoleID,
		FirstName: self.FirstName,
		LastName:  self.LastName,
		Email:     self.Email,
		Phone:     self.Phone,
		Password:  password,
		CreatedAt: self.CreatedAt,
		UpdatedAt: self.UpdatedAt,
	}
}

type JwtUser struct {
	UserID   string `json:"userId"`
	Email    string `json:"email"`
	RoleName string `json:"roleName"`
}

type UserContact struct {
	Phone string `json:"phone"`
}
