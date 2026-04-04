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
	UserId      string `json:"userId"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func New(title string, description string) *ReportPostRequest {
	return &ReportPostRequest{
		UserId:      "019d3c62-2f09-7456-ad9c-88813f3ecd13",
		Title:       title,
		Description: description,
	}

}

func PostReports(reportPostRequest *ReportPostRequest) {
	url := os.Getenv(AERITH_API_URL)
	payload, err := json.Marshal(reportPostRequest)
	body := bytes.NewBuffer(payload)

	response, err := http.Post(url+"/v1/reports", "application/json", body)
	if err != nil {
		log.Fatalf("%s", "Failed Aerith API POST "+err.Error())
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		log.Fatalf("Reddit Request nOK")
	}
}
