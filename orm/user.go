package orm

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string
	Password  string
	FullName  string
	AccountNo string
	Credit    int
}
