package data

import "gorm.io/gorm"

func Pagination(offset, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if offset <= 0 {
			offset = 0
		}

		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 20
		}

		offset = offset * limit
		return db.Offset(offset).Limit(limit)
	}
}
