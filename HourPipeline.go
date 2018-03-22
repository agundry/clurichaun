package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"./models"
)

const base_hour_hist_url = "https://min-api.cryptocompare.com/data/histohour"

/*
	Fetches crypto price data given url
 */
func cryptoGet(url string) (models.CryptoResponse, error) {

	var cryptoResponse models.CryptoResponse

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return cryptoResponse, err
	}

	// Construct client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return cryptoResponse, err
	}

	// Defer the closing of the body
	defer resp.Body.Close()

	// Read data into json format
	if err := json.NewDecoder(resp.Body).Decode(&cryptoResponse); err != nil {
		log.Println(err)
	}

	return cryptoResponse, nil
}

/*
	Runs crypto pipeline backfill process for each monitor for specified time range
 */
func cryptoPipeline(monitors []*models.Monitor, hoursBack int) {
	// Iterate over monitors
	for _, monitor := range monitors {
		// Format url
		url := fmt.Sprintf(base_hour_hist_url + "?fsym=%s&tsym=%s&limit=%d", monitor.Symbol, "USD", hoursBack)
		var cryptoResponse models.CryptoResponse

		// Fetch data
		cryptoResponse, err := cryptoGet(url)
		if err != nil {
			log.Fatal("Failed http get request: ", err)
		}

		// Process response
		models.InsertRecords(monitor.Id, cryptoResponse.Data)
	}
}

func main() {
	models.InitDB("root:my-secret-pw@tcp(localhost:13306)/clurichaun")
	monitors, err := models.FetchMonitors()
	if err != nil {
		log.Println(err)
	}

	cryptoPipeline(monitors, 60)
}
