package scrap

/*
This file is part of Alfred
(c) 2020 - 0xSha.io
*/

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/machinebox/graphql"
	"log"
	"time"
)

type H1response struct {

	HacktivityItems struct {
		Typename string `json:"__typename"`
		Edges    []struct {
			Typename string `json:"__typename"`
			Node     struct {
				Typename                    string    `json:"__typename"`
				DatabaseID                  string    `json:"databaseId"`
				ID                          string    `json:"id"`
				Reporter                    struct {
					Typename string `json:"__typename"`
					ID       string `json:"id"`
					Username string `json:"username"`
				} `json:"reporter"`

				Report struct {
					ID       string `json:"id"`
					Title    string `json:"title"`
					Substate string `json:"substate"`
					URL      string `json:"url"`
					Typename string `json:"__typename"`
				} `json:"report"`
				LatestDisclosableAction     string      `json:"latest_disclosable_action"`
				LatestDisclosableActivityAt time.Time   `json:"latest_disclosable_activity_at"`
				TotalAwardedAmount          interface{} `json:"total_awarded_amount"`
				SeverityRating              interface{} `json:"severity_rating"`
				Currency                    string      `json:"currency"`
				RequiresViewPrivilege bool `json:"requires_view_privilege"`
				Team                  struct {
					Typename             string `json:"__typename"`
					Handle               string `json:"handle"`
					ID                   string `json:"id"`
					MediumProfilePicture string `json:"medium_profile_picture"`
					Name                 string `json:"name"`
					URL                  string `json:"url"`
				} `json:"team"`
				Type               string      `json:"type"`
				Upvoted            interface{} `json:"upvoted"`
				Voters             struct {
					Typename string `json:"__typename"`
					Edges    []struct {
						Typename string `json:"__typename"`
						Node     struct {
							Typename string `json:"__typename"`
							ID       string `json:"id"`
							User     struct {
								Typename string `json:"__typename"`
								ID       string `json:"id"`
								Username string `json:"username"`
							} `json:"user"`
						} `json:"node"`
					} `json:"edges"`
				} `json:"voters"`
				Votes struct {
					Typename   string `json:"__typename"`
					TotalCount int    `json:"total_count"`
				} `json:"votes"`
			} `json:"node,omitempty"`
		} `json:"edges"`
		PageInfo struct {
			Typename    string `json:"__typename"`
			EndCursor   string `json:"endCursor"`
			HasNextPage bool   `json:"hasNextPage"`
		} `json:"pageInfo"`
		TotalCount int `json:"total_count"`
	} `json:"hacktivity_items"`
	Me interface{} `json:"me"`
}

func FetchHackerOne() (H1response, error) {

	client := graphql.NewClient("https://hackerone.com/graphql")

	// make a request
	req := graphql.NewRequest(`
		query HacktivityPageQuery($querystring: String, $orderBy: HacktivityItemOrderInput, $secureOrderBy: FiltersHacktivityItemFilterOrder, $where: FiltersHacktivityItemFilterInput, $maxShownVoters: Int) {
  me {
    id
    __typename
  }
  hacktivity_items(last: 25, after: "MjU", query: $querystring, order_by: $orderBy, secure_order_by: $secureOrderBy, where: $where) {
    total_count
    ...HacktivityList
    __typename
  }
}
fragment HacktivityList on HacktivityItemConnection {
  total_count
  pageInfo {
    endCursor
    hasNextPage
    __typename
  }
  edges {
    node {
      ... on HacktivityItemInterface {
        id
        databaseId: _id
        ...HacktivityItem
        __typename
      }
      __typename
    }
    __typename
  }
  __typename
}
fragment HacktivityItem on HacktivityItemUnion {
  type: __typename
  ... on HacktivityItemInterface {
    id
    votes {
      total_count
      __typename
    }
    voters: votes(last: $maxShownVoters) {
      edges {
        node {
          id
          user {
            id
            username
            __typename
          }
          __typename
        }
        __typename
      }
      __typename
    }
    upvoted: upvoted_by_current_user
    __typename
  }
  ... on Undisclosed {
    id
    ...HacktivityItemUndisclosed
    __typename
  }
  ... on Disclosed {
    id
    ...HacktivityItemDisclosed
    __typename
  }
  ... on HackerPublished {
    id
    ...HacktivityItemHackerPublished
    __typename
  }
}
fragment HacktivityItemUndisclosed on Undisclosed {
  id
  reporter {
    id
    username
    ...UserLinkWithMiniProfile
    __typename
  }
  team {
    handle
    name
    medium_profile_picture: profile_picture(size: medium)
    url
    id
    ...TeamLinkWithMiniProfile
    __typename
  }
  latest_disclosable_action
  latest_disclosable_activity_at
  requires_view_privilege
  total_awarded_amount
  currency
  __typename
}
fragment TeamLinkWithMiniProfile on Team {
  id
  handle
  name
  __typename
}
fragment UserLinkWithMiniProfile on User {
  id
  username
  __typename
}
fragment HacktivityItemDisclosed on Disclosed {
  id
  reporter {
    id
    username
    ...UserLinkWithMiniProfile
    __typename
  }
  team {
    handle
    name
    medium_profile_picture: profile_picture(size: medium)
    url
    id
    ...TeamLinkWithMiniProfile
    __typename
  }
  report {
    id
    title
    substate
    url
    __typename
  }
  latest_disclosable_action
  latest_disclosable_activity_at
  total_awarded_amount
  severity_rating
  currency
  __typename
}
fragment HacktivityItemHackerPublished on HackerPublished {
  id
  reporter {
    id
    username
    ...UserLinkWithMiniProfile
    __typename
  }
  team {
    id
    handle
    name
    medium_profile_picture: profile_picture(size: medium)
    url
    ...TeamLinkWithMiniProfile
    __typename
  }
  report {
    id
    url
    title
    substate
    __typename
  }
  latest_disclosable_activity_at
  severity_rating
  __typename
}
	`)

	ctx := context.Background()

	var resp H1response //  interface{}

	if err := client.Run(ctx, req, &resp); err != nil {
		log.Fatal(err)
		return resp,err
	}

	return resp, nil
}

func WriteH1ToDB(response H1response, entity Entity, db *gorm.DB) (int,error)  {

		totalFound := 0

		for _,report := range response.HacktivityItems.Edges{


			if report.Node.Report.Typename == "Report" {

				entity.Title = report.Node.Report.Title
				entity.URL = report.Node.Report.URL
				entity.Source = "HackerOne"

				if err := db.Create(&entity).Error; err !=nil {
					log.Println(err)
				}else {
					totalFound++
				}
				entity.ID++

				// ester egg
				// log.Println(report.Node.TotalAwardedAmount)

			}
		}
	return totalFound,nil

	}
