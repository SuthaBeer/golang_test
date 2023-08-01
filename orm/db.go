package orm

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB
var err error

func InitDB(database_name string) {
	Db, err := sql.Open("mysql", "root@tcp(localhost:3306)/")
	if err != nil {
		panic(err)
	}
	defer Db.Close()
	_, err = Db.Exec("CREATE DATABASE " + database_name)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = Db.Exec("USE " + database_name)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func InitDBColumns(database_name string) {
	host := os.Getenv("HOST")
	dsn := os.Getenv("MYSQL_DNS")
	dsn = strings.Replace(dsn, "HOST", host, 1)
	dsn = strings.Replace(dsn, "DATAGBASE_NAME", database_name, 1)

	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	Db.AutoMigrate(&User{})
	Db.AutoMigrate(&Transaction{})
}
