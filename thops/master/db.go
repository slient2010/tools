package master

import (
	_ "fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
)

var (
	Host     string
	Port     int
	Username string
	Password string
	Dbname   string
	Charset  string
)

func insertIntoDb(data *ServerInfo) {
	db, err := gorm.Open("mysql", Username+":"+Password+"@tcp("+Host+":"+strconv.Itoa(Port)+")/"+Dbname+"?charset="+Charset+"&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&ServerInfo{})

	db.Create(data)

}

func checkServerInOrNotInDb(data *ServerInfo) int {
	db, err := gorm.Open("mysql", Username+":"+Password+"@tcp("+Host+":"+strconv.Itoa(Port)+")/"+Dbname+"?charset="+Charset+"&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	count := 0
	db.Where(data).Find(data).Count(&count)
	return count

}
