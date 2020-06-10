package scrap

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"reflect"
	"testing"
)

func TestInitDB(t *testing.T) {

	got , err := InitDB()

	want := reflect.TypeOf(&gorm.DB{})

	if err != nil{
		t.Errorf("Init Error %s",err)
	}

	if reflect.TypeOf(got) != want{

		t.Errorf("Type Error %s",err)
	}

}