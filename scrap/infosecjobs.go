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
)

func FetchInfoSecJobs() (map[int][]string,error) {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.MaxDepth(1),

		colly.AllowedDomains("infosec-jobs.com"),
	)


	infoSecJobsArr := make(map[int][]string)

	// On every a element
	c.OnHTML("#job-list", func(e *colly.HTMLElement) {

		e.DOM.Find("a").Each(func(i int, s *goquery.Selection) {


			// only last 10 entries
			if i >= 10 {
				return
			}


			if s.Find("p").HasClass("job-list-item-company"){

				//fmt.Print(s.Find("p").First().Text() ) company

				link , _ := s.Attr("href")

				title := s.Find("p").Next().Text()

				infoSecJobsArr[i] = []string{title,link}

				//log.Print(link)

				//log.Println(title)
			}




		})


	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	// Start scraping
	err :=	c.Visit("https://infosec-jobs.com/")
   	if err !=nil {
		return infoSecJobsArr,err
	}

	return infoSecJobsArr,nil
}


func WriteInfoSecJobsToDB(newsArr map[int][]string,entity Entity, db *gorm.DB) (int,error)  {

	totalFound := 0

	for _,key := range newsArr{

		entity.Title = key[0]
		entity.URL = "https://infosec-jobs.com/"+key[1]
		entity.Source = "InfoSecJobs"

		if err := db.Create(&entity).Error; err !=nil {
			log.Println(err)
		}else {
			totalFound++
		}
		entity.ID++
	}
	return totalFound,nil

}

