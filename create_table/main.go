package main

import (
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
)

type Record struct {
    gorm.Model
}

func getDB() (*gorm.DB, error){
    user := "root"
    passwd := "123456"
    host := "localhost"
    port := "33070"
    dbname := "large_table_for_golang"
    url := user + ":" + passwd + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8&parseTime=true"
    db, err := gorm.Open("mysql", url)
    return db, err
}

func main() {
    db, err := getDB()
    if err!=nil{
        return
    }
    defer db.Close()
    db.AutoMigrate(&Record{})
    count := 1000000
    for i := 0; i < count; i++ {
        if i%1000 == 0 {
            fmt.Println(i, "/", count)
        }
        db.Create(&Record{})
    }

}
