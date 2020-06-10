package scrap

/*
This file is part of Alfred
(c) 2020 - 0xSha.io
*/

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	"log"
	"strings"
)




func FetchGithubAdvisory()  (map[int][]string , error) {
	c := colly.NewCollector(
		colly.MaxDepth(1),

		colly.AllowedDomains("github.com"),
	)


	githubAdvisoryArr := make(map[int][]string)


	// On every a element
	c.OnHTML(".Box", func(e *colly.HTMLElement) {

		e.DOM.Find(".Box-row").Each(func(i int, s *goquery.Selection) {

			// only last 10 entries
			if i >= 10 {
				return
			}



			link , _ := s.First().Find("a").Attr("href")

			title := strings.TrimSpace(s.First().Find("a").Text())


			//log.Println(title)

			githubAdvisoryArr[i] = []string{link,title}
			log.Println(title)

			//hackerNewsArr = append(arr["URL"], &hackerNewsArr )




		})


	})

	//fmt.Println(hackerNewsArr)

	err := c.Visit("https://github.com/advisories")
	if err !=nil {
		return githubAdvisoryArr,err
	}

	return githubAdvisoryArr,nil

}


func WriteGithubAdvisoryToDB(newsArr map[int][]string,entity Entity, db *gorm.DB) (int,error)  {

	totalFound := 0

	for _,key := range newsArr{

		entity.Title = key[0]
		entity.URL = "https://github.com/advisories"+key[1]
		entity.Source = "GitHubAdvisory"

		if err := db.Create(&entity).Error; err !=nil {
			log.Println(err)
		}else {
			totalFound++
		}
		entity.ID++
	}
	return totalFound,nil

}

