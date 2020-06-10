package scrap

/*
This file is part of Alfred
(c) 2020 - 0xSha.io
*/

import (
	"github.com/jinzhu/gorm"
)



type Entity struct {
	gorm.Model
	Title string `gorm:"type:varchar(255);unique"`
	URL string `gorm:"type:varchar(255);unique"`
	Source string `gorm:"type:varchar(255)"`
}

type Master struct {
	gorm.Model
	Name string `gorm:"type:varchar(255);unique"`
}


func InitDB() ( *gorm.DB , error) {


	db, err := gorm.Open("sqlite3", "pennyworth.db")
	if err != nil {
		return nil,err
	}


	db.AutoMigrate(&Entity{})
	db.AutoMigrate(&Master{})

	return db,nil

}
