package scrap

/*
This file is part of Alfred
(c) 2020 - 0xSha.io
*/

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"time"
)



const endPoint = "https://newsapi.org/v2/everything/"
// same as other modules only 10 results
const pageSize =  10
// pennyworth API key
const key = "a780abdb06b64ffdbcf8223612a8fbf6"


type NewsApiResp struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			ID   interface{} `json:"id"`
			Name string      `json:"name"`
		} `json:"source"`
		Author      string    `json:"author"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		URL         string    `json:"url"`
		URLToImage  string    `json:"urlToImage"`
		PublishedAt time.Time `json:"publishedAt"`
		Content     string    `json:"content"`
	} `json:"articles"`
}



func FetchNewsAPI() (map[int][]string,error) {
	// selected stuff
	keywords := []string{"0-day", "hacker",  "data-breach" , "bug-bounty" ,"vulnerability" , "malware"}

	newsAPIArr := make(map[int][]string)

	counter := 0

	var newsApi NewsApiResp
	for i := 0; i<len(keywords);i++{
		query :=  fmt.Sprintf("?qInTitle=%s&pagesize=%d&sortBy=publishedAt&language=en&apiKey=%s" , keywords[i] , pageSize, key)

		resp , err := http.Get(endPoint +query)

		if err != nil {
			return nil,err
			log.Println(err)
		}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&newsApi)

		if err != nil {

			return nil,err

			log.Println(err)
		}

		counter++

		for _, article := range newsApi.Articles{

			title := article.Title
			url := article.URL

			newsAPIArr[counter] = []string{title,url}
			log.Println(article.Title)
			log.Println(article.URL)

			counter++
		}


	}

	return newsAPIArr,nil
}

func WriteNewsAPIToDB(newsArr map[int][]string,entity Entity, db *gorm.DB) (int,error)  {

	totalFound := 0

	for _,key := range newsArr{

		entity.Title = key[0]
		entity.URL = key[1]
		entity.Source = "NewsAPI"

		if err := db.Create(&entity).Error; err !=nil {
			log.Println(err)
		}else {
			totalFound++
		}
		entity.ID++
	}
	return totalFound,nil
}
