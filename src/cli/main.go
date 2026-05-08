package main

import (
	"cloud-strife-worker/services"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type IngestionCheckpoint struct {
	TotalProcessed int                  `json:"totalProcessed"`
	After          *string              `json:"after"`
	Reason         string               `json:"reason"`
	LastReport     *services.RedditPost `json:"lastReport"`
	CreatedAt      string               `json:"createdAt"`
}

func saveCheckpoint(
	total int,
	after *string,
	reason string,
	lastReport *services.RedditPost,
) {
	checkpoint := IngestionCheckpoint{
		TotalProcessed: total,
		After:          after,
		Reason:         reason,
		LastReport:     lastReport,
		CreatedAt:      time.Now().Format(time.RFC3339),
	}

	file, err := os.Create("checkpoint.json")
	if err != nil {
		log.Printf("failed to create checkpoint: %v", err)
		return
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(checkpoint); err != nil {
		log.Printf("failed to write checkpoint: %v", err)
	}
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var after *string = nil
	maxPosts := 500
	total := 0

	var lastReport *services.RedditPost = nil

	for {
		getRedditDto := services.GetPostBydSubRedditRequest{
			Subreddit: "golpes",
			Limit:     5,
			After:     after,
		}

		posts, nextAfter := services.GetPostsBySubReddit(&getRedditDto)

		if len(posts) == 0 {
			log.Println("No posts found")

			saveCheckpoint(
				total,
				after,
				"NO_POSTS_FOUND",
				lastReport,
			)

			break
		}

		for _, report := range posts {
			lastReport = &report

			reportReq := services.New(
				report,
				getRedditDto.Subreddit,
			)

			services.PostReports(reportReq)

			total++
		}

		if total >= maxPosts {
			log.Println("Reached maxPosts")

			saveCheckpoint(
				total,
				after,
				"MAX_POSTS_REACHED",
				lastReport,
			)

			break
		}

		if nextAfter == nil {
			log.Println("No more pages")

			saveCheckpoint(
				total,
				after,
				"NO_MORE_PAGES",
				lastReport,
			)

			break
		}

		after = nextAfter

		time.Sleep(2 * time.Second)
	}
}
