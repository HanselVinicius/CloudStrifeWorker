package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

const AERITH_API_URL = "AERITH_API_URL"

type ReportPostRequest struct {
	UserId        string  `json:"userId"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	RedditID      string  `json:"redditId"`
	Score         int     `json:"score"`
	Ups           int     `json:"ups"`
	IsVideo       bool    `json:"isVideo"`
	URL           string  `json:"url"`
	MediaType     *string `json:"mediaType"`
	SubredditName string  `json:"subredditName"`
}

func New(post RedditPost, subreddit string) *ReportPostRequest {
	var mediaType *string = nil

	if post.Media != nil {
		mediaType = &post.Media.Type
	}

	return &ReportPostRequest{
		UserId:        "019d3c62-2f09-7456-ad9c-88813f3ecd13",
		Title:         post.Title,
		Description:   post.SelfText,
		RedditID:      post.ID,
		Score:         post.Score,
		Ups:           post.Ups,
		IsVideo:       post.IsVideo,
		URL:           post.URLOverridden,
		MediaType:     mediaType,
		SubredditName: subreddit,
	}
}

func PostReports(reportPostRequest *ReportPostRequest) {
	url := os.Getenv(AERITH_API_URL)

	payload, err := json.Marshal(reportPostRequest)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	body := bytes.NewBuffer(payload)

	response, err := http.Post(
		url+"/v1/reports",
		"application/json",
		body,
	)

	if err != nil {
		log.Printf("%s", "Failed Aerith API POST "+err.Error())
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		log.Printf("Aerith API Request nOK")
	}
}
