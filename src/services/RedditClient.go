package services

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

const REDDIT_API_URL = "REDDIT_API_URL"

type GetPostBydSubRedditRequest struct {
	Subreddit string
	Limit     int
	Offset    int
}

type RedditMedia struct {
	Type string `json:"type"`
}

type RedditPost struct {
	ID       string       `json:"id"`
	Title    string       `json:"title"`
	SelfText string       `json:"selftext"`
	IsVideo  bool         `json:"is_video"`
	Ups      int          `json:"ups"`
	Score    int          `json:"score"`
	Media    *RedditMedia `json:"media"`
}

type RedditListing struct {
	Data struct {
		Children []struct {
			Data RedditPost `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func GetPostsBySubReddit(getPostBydSubRedditRequest *GetPostBydSubRedditRequest) []RedditPost {
	url := os.Getenv(REDDIT_API_URL)

	response, err := http.Get(url + "r/" + getPostBydSubRedditRequest.Subreddit)
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

	return posts
}
