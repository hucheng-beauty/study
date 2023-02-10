package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string    `gorm:"primarykey;type:uuid;comment:主键ID"`
	CreatedAt time.Time `gorm:"column:created_at;sort:desc;index:idx_created_at;comment:创建时间"`
	UpdatedAt time.Time `gorm:"column:updated_at;comment:修改时间"`
	DeletedAt gorm.DeletedAt
}

type StrSlice []string

func (j *StrSlice) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), j)
}

func (j StrSlice) Value() (driver.Value, error) {
	if j == nil {
		return json.Marshal(StrSlice{})
	}
	return json.Marshal(j)
}
