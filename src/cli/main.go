package main

import (
	"cloud-strife-worker/services"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	getRedditDto := services.GetPostBydSubRedditRequest{
		Subreddit: "golpes",
		Limit:     5,
		Offset:    0,
	}

	body := services.GetPostsBySubReddit(&getRedditDto)

	log.Print(body)

}
