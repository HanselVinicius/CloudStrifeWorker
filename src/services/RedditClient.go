package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const REDDIT_API_URL = "REDDIT_API_URL"

type GetPostBydSubRedditRequest struct {
	Subreddit string
	Limit     int
	After     *string
}

type RedditMedia struct {
	Type string `json:"type"`
}

type RedditPost struct {
	ID            string       `json:"id"`
	Title         string       `json:"title"`
	SelfText      string       `json:"selftext"`
	IsVideo       bool         `json:"is_video"`
	Ups           int          `json:"ups"`
	Score         int          `json:"score"`
	Media         *RedditMedia `json:"media"`
	URLOverridden string       `json:"url_overridden_by_dest"`
}

type RedditListing struct {
	Data struct {
		After    string `json:"after"`
		Children []struct {
			Data RedditPost `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func GetPostsBySubReddit(getPostBydSubRedditRequest *GetPostBydSubRedditRequest) ([]RedditPost, *string) {
	url := os.Getenv(REDDIT_API_URL)

	reqUrl := url + "r/" + getPostBydSubRedditRequest.Subreddit + "?limit="

	reqUrl += fmt.Sprintf("%d", getPostBydSubRedditRequest.Limit)

	if getPostBydSubRedditRequest.After != nil {
		reqUrl += "&after=" + *getPostBydSubRedditRequest.After
	}

	response, err := http.Get(reqUrl)
	if err != nil {
		log.Fatalf("Failed Reddit Request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Reddit Request nOK")
	}

	var listing RedditListing
	json.NewDecoder(response.Body).Decode(&listing)

	posts := make([]RedditPost, 0, len(listing.Data.Children))
	for _, child := range listing.Data.Children {
		posts = append(posts, child.Data)
	}

	if listing.Data.After == "" {
		return posts, nil
	}

	return posts, &listing.Data.After
}
