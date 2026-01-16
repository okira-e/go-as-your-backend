package models

import "github.com/samborkent/uuidv7"

type Role struct {
	ID   string `sql:"id" gorm:"type:uuid;primaryKey"`
	Name string `sql:"name" gorm:"size:32;uniqueIndex;not null"`
}

func (self *Role) ToDto() *RoleDto {
	return &RoleDto{
		ID:   self.ID,
		Name: self.Name,
	}
}

type CreateRoleDto struct {
	Name string `json:"name" validate:"required"`
}

func (self *CreateRoleDto) FromDto() *Role {
	id := uuidv7.New().String()

	return &Role{
		ID:   id,
		Name: self.Name,
	}
}

type RoleDto struct {
	ID   string `json:"id" `
	Name string `json:"name" validate:"required"`
}
