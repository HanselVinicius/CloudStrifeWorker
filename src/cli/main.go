package main

import (
	"cloud-strife-worker/services"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var after *string = nil
	maxPosts := 50
	total := 0

	for {
		getRedditDto := services.GetPostBydSubRedditRequest{
			Subreddit: "golpes",
			Limit:     5,
			After:     after,
		}

		posts, nextAfter := services.GetPostsBySubReddit(&getRedditDto)

		if len(posts) == 0 {
			log.Println("No posts found")
			break
		}

		for _, report := range posts {
			reportReq := services.New(report.Title, report.SelfText)
			services.PostReports(reportReq)
			total++
		}

		if total >= maxPosts {
			log.Println("Reached maxPosts")
			break
		}

		if nextAfter == nil {
			log.Println("No more pages")
			break
		}

		after = nextAfter

		time.Sleep(2 * time.Second)
	}
}
