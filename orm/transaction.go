package orm

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserID   uint
	ToUserID uint
	Credit   int
}
