package main

import (
	"crawler/engine"
	"crawler/scheduler"
	"crawler/wiki/parser"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
)
import "crawler/types"

type Item struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

func main() {
	contents, err := ioutil.ReadFile("./seed.conf")
	if err != nil {
		log.Fatal(err)
	}

	var items []Item
	err = json.Unmarshal(contents, &items)
	if err != nil {
		log.Fatal(err)
	}

	requests := []types.Request{}

	for _, item := range items {
		requests = append(requests, types.Request{
			Url:       item.Url,
			Title:     item.Title,
			ParseFunc: parser.ParseWiki,
		})
	}

	fmt.Println(requests)

	runtime.GOMAXPROCS(runtime.NumCPU())

	e := engine.ConcurrentEngine{
		Scheduler:   &scheduler.QueuedScheduler{},
		WorkerCount: 100,
	}

	//e := engine.SimpleEngine{}
	e.Run(requests...)

	//a := []byte("abcd")
	//b := []byte("efg")
	//fmt.Println(string(append(a, b...)))

	//engine, err := xorm.NewEngine("mysql", "root:aa0987aa@(localhost:3306)/wiki?charset=utf8")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//engine.ShowSQL(true)
	//
	//engine.SetMaxOpenConns(10)
	//engine.SetMaxIdleConns(10)
	//engine.SetConnMaxLifetime(3600 * time.Second)
	//
	//content := []engine2.TitleUrl{}
	//
	//err = engine.Where("title=?", "bitcoin1").Find(&content)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Println(content)
}
