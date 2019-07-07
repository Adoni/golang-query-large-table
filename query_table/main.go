package main

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/go-sql-driver/mysql"
    "time"
)

type Record struct {
    gorm.Model
    Id uint `gorm:"column:object_id"`
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

func timeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    fmt.Printf("%s tooks %s\n", name, elapsed)
}

func method1(db *gorm.DB) {
    records := make([]Record, 0)
    defer timeTrack(time.Now(), "Method 1")
    db.Table("records").Select("id").Limit(1000000).Find(&records)
}

func method2(db *gorm.DB) {
    records := make([]Record, 0)
    defer timeTrack(time.Now(), "Method 2")
    count := 1000000
    bucketSize := 100000
    for beginId :=1; beginId <=count; beginId += bucketSize {
        fmt.Printf("%d of %d\n", beginId, count)
        currentRecords := make([]Record, 0)
        db.Table("records").Select("id").Where("id>=? and id<?", beginId, beginId+bucketSize).Find(&currentRecords)
        records = append(records, currentRecords...)
    }
}

func method3(db *gorm.DB) {
    records := make([]Record, 0)
    defer timeTrack(time.Now(), "Method 3")
    count := 1000000
    bucketSize := 100000
    resultCount := 0
    resultChannel := make(chan []Record, 0)
    for beginID :=1; beginID <=count; beginID += bucketSize {
        endId := beginID +bucketSize
        go func(beginId int, endId int) {
            currentRecords := make([]Record, 0)
            db.Table("records").Select("id").Where("id>=? and id<?", beginId, endId).Find(&currentRecords)
            resultChannel <- currentRecords
        }(beginID, endId)
        resultCount += 1
    }
    for i:=0; i<resultCount; i++{
        fmt.Printf("%d of %d\n", i, resultCount)
        currentRecords := <- resultChannel
        records = append(records, currentRecords...)
    }
}


func main(){
    db, err := getDB()
    if err!=nil{
        fmt.Println(err)
        return
    }
    defer db.Close()
    method1(db)
    method2(db)
    method3(db)
}
