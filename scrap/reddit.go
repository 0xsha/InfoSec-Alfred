package scrap

/*
This file is part of Alfred
(c) 2020 - 0xSha.io
*/

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"html"
	"log"
	"net/http"
	"strings"
)

type RedditJson struct {
	Kind string `json:"kind"`
	Data struct {
		Modhash  string `json:"modhash"`
		Dist     int    `json:"dist"`
		Children []struct {
			Kind string `json:"kind"`
			Data struct {
				ApprovedAtUtc              interface{}   `json:"approved_at_utc"`
				Subreddit                  string        `json:"subreddit"`
				Selftext                   string        `json:"selftext"`
				AuthorFullname             string        `json:"author_fullname"`
				Saved                      bool          `json:"saved"`
				ModReasonTitle             interface{}   `json:"mod_reason_title"`
				Gilded                     int           `json:"gilded"`
				Clicked                    bool          `json:"clicked"`
				Title                      string        `json:"title"`
				LinkFlairRichtext          []interface{} `json:"link_flair_richtext"`
				SubredditNamePrefixed      string        `json:"subreddit_name_prefixed"`
				Hidden                     bool          `json:"hidden"`
				Pwls                       int           `json:"pwls"`
				LinkFlairCSSClass          interface{}   `json:"link_flair_css_class"`
				Downs                      int           `json:"downs"`
				TopAwardedType             interface{}   `json:"top_awarded_type"`
				HideScore                  bool          `json:"hide_score"`
				Name                       string        `json:"name"`
				Quarantine                 bool          `json:"quarantine"`
				LinkFlairTextColor         string        `json:"link_flair_text_color"`
				UpvoteRatio                float64       `json:"upvote_ratio"`
				AuthorFlairBackgroundColor interface{}   `json:"author_flair_background_color"`
				SubredditType              string        `json:"subreddit_type"`
				Ups                        int           `json:"ups"`
				TotalAwardsReceived        int           `json:"total_awards_received"`
				Created                  float64       `json:"created"`
				LinkFlairType            string        `json:"link_flair_type"`
				Wls                      int           `json:"wls"`
				RemovedByCategory        interface{}   `json:"removed_by_category"`
				BannedBy                 interface{}   `json:"banned_by"`
				AuthorFlairType          string        `json:"author_flair_type"`
				Domain                   string        `json:"domain"`
				AllowLiveComments        bool          `json:"allow_live_comments"`
				SelftextHTML             interface{}   `json:"selftext_html"`
				Likes                    interface{}   `json:"likes"`
				SuggestedSort            interface{}   `json:"suggested_sort"`
				BannedAtUtc              interface{}   `json:"banned_at_utc"`
				ViewCount                interface{}   `json:"view_count"`
				Archived                 bool          `json:"archived"`
				NoFollow                 bool          `json:"no_follow"`
				IsCrosspostable          bool          `json:"is_crosspostable"`
				Pinned                   bool          `json:"pinned"`
				Over18                   bool          `json:"over_18"`
				AllAwardings             []interface{} `json:"all_awardings"`
				Awarders                 []interface{} `json:"awarders"`
				MediaOnly                bool          `json:"media_only"`
				CanGild                  bool          `json:"can_gild"`
				Spoiler                  bool          `json:"spoiler"`
				Locked                   bool          `json:"locked"`
				AuthorFlairText          interface{}   `json:"author_flair_text"`
				TreatmentTags            []interface{} `json:"treatment_tags"`
				Visited                  bool          `json:"visited"`
				RemovedBy                interface{}   `json:"removed_by"`
				NumReports               interface{}   `json:"num_reports"`
				Distinguished            interface{}   `json:"distinguished"`
				SubredditID              string        `json:"subreddit_id"`
				ModReasonBy              interface{}   `json:"mod_reason_by"`
				RemovalReason            interface{}   `json:"removal_reason"`
				LinkFlairBackgroundColor string        `json:"link_flair_background_color"`
				ID                       string        `json:"id"`
				IsRobotIndexable         bool          `json:"is_robot_indexable"`
				ReportReasons            interface{}   `json:"report_reasons"`
				Author                   string        `json:"author"`
				DiscussionType           interface{}   `json:"discussion_type"`
				NumComments              int           `json:"num_comments"`
				SendReplies              bool          `json:"send_replies"`
				WhitelistStatus          string        `json:"whitelist_status"`
				ContestMode              bool          `json:"contest_mode"`
				ModReports               []interface{} `json:"mod_reports"`
				AuthorPatreonFlair       bool          `json:"author_patreon_flair"`
				CrosspostParent          string        `json:"crosspost_parent"`
				AuthorFlairTextColor     interface{}   `json:"author_flair_text_color"`
				Permalink                string        `json:"permalink"`
				ParentWhitelistStatus    string        `json:"parent_whitelist_status"`
				Stickied                 bool          `json:"stickied"`
				URL                      string        `json:"url"`
				SubredditSubscribers     int           `json:"subreddit_subscribers"`
				CreatedUtc               float64       `json:"created_utc"`
				NumCrossposts            int           `json:"num_crossposts"`
				Media                    interface{}   `json:"media"`
				IsVideo                  bool          `json:"is_video"`
			} `json:"data"`
		} `json:"children"`
		After  interface{} `json:"after"`
		Before string      `json:"before"`
	} `json:"data"`
}

func FetchNetSecReddit() (map[int][]string, error) {


	client := &http.Client{}


	// Didn't know such endpoint exists
	req, err := http.NewRequest("GET", "https://www.reddit.com/r/netsec/.json?count=10", nil)
	if err != nil {
		return nil, err
		log.Println(err)
	}

	// well we also need a bypass for reddit client check
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {

		return nil,err
		log.Println(err)
	}

	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	//var data []byte

	redditNewsArr := make(map[int][]string)


	var redditResp RedditJson

	err = json.NewDecoder(resp.Body).Decode(&redditResp)

	if err != nil {
		log.Println(err)
		return nil,err
	}

	//fmt.Println(redditResp.Data)

	for i,item := range redditResp.Data.Children {

		if i >= 10 {
			break
		}

		title := item.Data.Title
		url := item.Data.URL

		title = strings.TrimSpace(html.UnescapeString(title))

		redditNewsArr[i] = []string{title,url}


		log.Println(title)
		log.Println(url)

	}

	return redditNewsArr,nil

	//log.Println(string(body))



}


func WriteRedditNetSecToDB(newsArr map[int][]string,entity Entity, db *gorm.DB) (int,error)  {

	totalFound := 0

	for _,key := range newsArr{

		entity.Title = key[0]
		entity.URL = key[1]
		entity.Source = "Reddit/NetSEC"

		if err := db.Create(&entity).Error; err !=nil {
			log.Println(err)
		}else {
			totalFound++
		}
		entity.ID++
	}
	return totalFound,nil

}
