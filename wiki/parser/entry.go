package parser

import (
	"crawler/engine"
	"crawler/types"
	"log"
	"regexp"
)

//const entryRe = `<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`
//const entryRe = `<a href="(https://zh.wikipedia.org/wiki/\S+)" title="\S+">(\S+)</a>`
const entryRe = `<a href="(/wiki/\S+)" title="\S+">(\S+)</a>`
const urlDomain = "https://zh.wikipedia.org"

var re *regexp.Regexp

func init() {
	var err error
	re, err = regexp.Compile(entryRe)
	if err != nil {
		log.Fatal(err)
	}
}

func ParseEntry(contents []byte) types.ParseResult {
	matches := re.FindAllSubmatch(contents, -1)
	result := types.ParseResult{}

	for _, m := range matches {
		title := string(m[2])
		url := string(m[1])

		titleUrl := []engine.TitleUrl{}
		xormEngine := engine.GetEngine()
		err := xormEngine.Where("url=?", url).Find(&titleUrl)
		if err != nil {
			//fmt.Printf("Find url: %s in database error: %s\n", url, err)
			return types.ParseResult{}
		}

		if len(titleUrl) >= 1 {
			//fmt.Printf("Find url: %s in database has exist\n", url)
			return types.ParseResult{}
		}

		result.Items = append(result.Items, "Title: "+title)
		result.Requests = append(
			result.Requests,
			types.Request{
				Url:   url,
				Title: title,
				ParseFunc: func(bytes []byte) types.ParseResult {
					return ParseEntry(bytes)
				},
			})
	}
	return result
}

func ParseWiki(contents []byte) types.ParseResult {
	matches := re.FindAllSubmatch(contents, -1)
	result := types.ParseResult{}

	for _, m := range matches {
		title := string(m[2])
		url := urlDomain + string(m[1])

		content := []engine.Content{}
		xormEngine := engine.GetEngine()
		err := xormEngine.Where("downloadurl=?", url).Find(&content)
		if err != nil {
			//fmt.Printf("Find downloadurl: %s in database error: %s\n", url, err)
			continue
		}

		if len(content) >= 1 {
			//fmt.Printf("Find title: %s in database has exist\n", title)
			continue
		}

		result.Items = append(result.Items, "Title: "+title)
		result.Requests = append(
			result.Requests,
			types.Request{
				Url:   url,
				Title: title,
				ParseFunc: func(bytes []byte) types.ParseResult {
					return ParseWiki(bytes)
				},
			})
	}
	return result
}
