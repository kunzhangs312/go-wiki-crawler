package engine

import (
	"crawler/fetcher"
	"crawler/types"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"time"
)

type SimpleEngine struct{}

// 串行
func (e SimpleEngine) Run(seeds ...types.Request) {
	var requests []types.Request
	for _, r := range seeds {
		requests = append(requests, r)
	}

	for len(requests) > 0 {
		r := requests[0]
		requests = requests[1:]

		parseResult, err := worker(r)
		if err != nil {
			continue
		}

		requests = append(requests, parseResult.Requests...)
		for _, item := range parseResult.Items {
			log.Printf("Got item %v", item)
		}
	}
}

//func worker(r types.Request) (types.ParseResult, error) {
//	log.Printf("Fetching %s", r.Url)
//	body, err := fetcher.Fetch(r.Url)
//	if err != nil {
//		log.Printf("Fetcher: error fetching url %s: %v", r.Url, err)
//		return types.ParseResult{}, err
//	}
//	return r.ParseFunc(body), nil
//}

type TitleUrl struct {
	Id    int    `xorm:"not null pk autoincr INT(11)"`
	Title string `xorm:"not null VARCHAR(255)"`
	Url   string `xorm:"not null unique VARCHAR(1024)"`
	Hash  string `xorm: VARCHAR(255)"`
}

type Content struct {
	Platform            string    `xorm:"VARCHAR(255)"`
	Channelname         string    `xorm:"VARCHAR(255)"`
	Channel             string    `xorm:"VARCHAR(255)"`
	Label               string    `xorm:"VARCHAR(255)"`
	Describe            string    `xorm:"VARCHAR(2048)"`
	Downloadurl         string    `xorm:"unique VARCHAR(1024)"`
	Category            string    `xorm:"VARCHAR(255)"`
	Size                int       `xorm:"INT(11)"`
	Addtime             time.Time `xorm:"DATETIME"`
	Id                  int       `xorm:"not null pk autoincr INT(11)"`
	Poster              string    `xorm:"VARCHAR(255)"`
	Extra               string    `xorm:"VARCHAR(255)"`
	IsMainland          int       `xorm:"not null comment('大陆是否可以下载') TINYINT(1)"`
	IssuedStatus        string    `xorm:"not null default 'unissued' comment('下发状态') index VARCHAR(255)"`
	LastIssuedTime      time.Time `xorm:"DATETIME"`
	MineStatus          string    `xorm:"not null default 'unminer' comment('挖矿状态') index VARCHAR(255)"`
	IssuedMachineNumber int       `xorm:"not null default 00000000000 comment('下发矿机个数') INT(11)"`
	MineMachineNumber   int       `xorm:"not null default 00000000000 comment('参与挖矿的矿机个数，这个值只是统计上报挖矿的主机个数，不去重') INT(11)"`
}


var xormEngine *xorm.Engine

func init() {
	engine, err := xorm.NewEngine("mysql", "root:123456@(192.168.0.63:3306)/poseidon_content?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}

	engine.ShowSQL(false)

	engine.SetMaxOpenConns(10)
	engine.SetMaxIdleConns(10)
	engine.SetConnMaxLifetime(3600 * time.Second)

	xormEngine = engine
}

func GetEngine() *xorm.Engine {
	return xormEngine
}

func workerWiki(r types.Request) (types.ParseResult, error) {
	content := []Content{}
	err := xormEngine.Where("downloadurl=?", r.Url).Find(&content)
	if err != nil {
		//fmt.Printf("Find downloadurl: %s in database error: %s\n", r.Url, err)
		return types.ParseResult{}, err
	}

	if len(content) >= 1 {
		//fmt.Printf("Find title: %s url: %s in database has exist\n", r.Title, r.Url)
		return types.ParseResult{}, fmt.Errorf("Find downloadurl: %s in database has exist\n", r.Url)
	}

	contentEntry := Content{Platform:"wikipedia", Channelname:"None", Channel:"None", Label:r.Title, Describe:"None",
		Downloadurl:r.Url, Category:"html", Size:-1, Addtime:time.Now(), Poster:"None", IsMainland:0,
		IssuedStatus:"unissued", MineStatus:"unminer", IssuedMachineNumber:0, MineMachineNumber:0}

	_, err = xormEngine.Insert(&contentEntry)
	if err != nil {
		fmt.Printf("Insert database url: %s error: %s\n", r.Url, err)
	}

	log.Printf("Fetching title: %s, url: %s", r.Title, r.Url)
	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		log.Printf("Fetcher: error fetching url %s: %v", r.Url, err)
		return types.ParseResult{}, err
	}

	//fmt.Printf("Insert database url: %s success, num: %d\n", r.Url, num)

	//fmt.Println(string(body))

	return r.ParseFunc(body), nil
}

func worker(r types.Request) (types.ParseResult, error) {
	log.Printf("Fetching Title: %s, Url: %s", r.Title, r.Url)

	titleUrl := []TitleUrl{}
	err := xormEngine.Where("url=?", r.Url).Find(&titleUrl)
	if err != nil {
		fmt.Printf("Find url: %s in database error: %s\n", r.Url, err)
		return types.ParseResult{}, err
	}

	if len(titleUrl) >= 1 {
		fmt.Printf("Find url: %s in database has exist\n", r.Url)
		return types.ParseResult{}, fmt.Errorf("Find url: %s in database has exist\n", r.Url)
	}

	fr := fetcher.NewDefaultWikiDownloader()
	fr.SetDownloadNo(1)
	fr.SetTmpDownloadDir(r.Title)
	body, err := fr.Download(r.Url, r.Title)
	if err != nil {
		fmt.Printf("Url: %s fetch error: %s\n", r.Url, err)
		return types.ParseResult{}, err
	}

	titleUrlEntry := TitleUrl{Url:r.Url, Title:r.Title}

	num, err := xormEngine.Insert(&titleUrlEntry)
	if err != nil {
		fmt.Printf("Insert database url: %s error: %s\n", r.Url, err)
	}

	fmt.Printf("Insert database url: %s success, num: %d\n", r.Url, num)

	return r.ParseFunc(body), nil
}
