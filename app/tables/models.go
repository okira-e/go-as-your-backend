package tables

import (
	"database/sql"
	"time"
)

// SystemUsers is the model for the system_users table
type systemUsers struct {
	ID        uint `gorm:"primaryKey, autoIncrement"`
	Name      string
	Email     *string         `gorm:"unique"`
	Salary    sql.NullFloat64 `gorm:"type:decimal(10,2)"`
	Age       sql.NullInt32
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

var Models = []any{
	&systemUsers{},
}
