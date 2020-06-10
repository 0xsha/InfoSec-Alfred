package scrap

/*
This file is part of Alfred
(c) 2020 - 0xSha.io
*/

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/jinzhu/gorm"
	"log"
	"strings"
)




func FetchTheHackerNews()  (map[int][]string , error) {
	c := colly.NewCollector(
		colly.MaxDepth(1),

		colly.AllowedDomains("thehackernews.com"),
	)


	hackerNewsArr := make(map[int][]string)


	// On every a element
	c.OnHTML("#Blog1", func(e *colly.HTMLElement) {

		e.DOM.Find(".story-link").Each(func(i int, s *goquery.Selection) {

			// only last 10 entries
			if i >= 10 {
				return
			}


			link , _ := s.Attr("href")

			title := strings.TrimSpace(s.Find(".home-title").Text())


			//log.Println(title)

			hackerNewsArr[i] = []string{title,link}
			//log.Println(link)

			//hackerNewsArr = append(arr["URL"], &hackerNewsArr )




		})


	})

	//fmt.Println(hackerNewsArr)

	err := c.Visit("https://thehackernews.com/search?max-results=10")
	if err !=nil {
		return hackerNewsArr,err
	}

	return hackerNewsArr,nil

}

func WriteHNewsToDB(newsArr map[int][]string,entity Entity, db *gorm.DB) (int,error)  {

	totalFound := 0

	for _,key := range newsArr{

		entity.Title = key[0]
		entity.URL = key[1]
		entity.Source = "TheHackerNews"

		if err := db.Create(&entity).Error; err !=nil {
			log.Println(err)
		}else {
			totalFound++
		}
		entity.ID++
	}
	return totalFound,nil

}
