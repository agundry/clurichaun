package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "net/url"
	_ "github.com/go-sql-driver/mysql"
)

const base_hour_hist_url = "https://min-api.cryptocompare.com/data/histohour"

type HourHist struct {
	Time       int     `json:"time"`
	Close      float64 `json:"close"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Open       float64 `json:"open"`
	Volumefrom float64 `json:"volumefrom"`
	Volumeto   float64 `json:"volumeto"`
}

type CryptoResponse struct {
	Response   string `json:"Response"`
	Type       int    `json:"Type"`
	Aggregated bool   `json:"Aggregated"`
	Data       []HourHist `json:"Data"`
	TimeTo            int  `json:"TimeTo"`
	TimeFrom          int  `json:"TimeFrom"`
	FirstValueInArray bool `json:"FirstValueInArray"`
	ConversionType    struct {
		Type             string `json:"type"`
		ConversionSymbol string `json:"conversionSymbol"`
	} `json:"ConversionType"`
}

func insertRecords(data []HourHist) {
	db, err := sql.Open("mysql", "root:my-secret-pw@tcp(localhost:13306)/clurichaun")
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Construct prepared statement
	stmtIns, err := db.Prepare("INSERT INTO `records` (monitor_id, record_epoch, value, updated) VALUES(?,?,?,NOW())") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Insert price records into db
	for i := range data {
		item := data[i]
		_, err = stmtIns.Exec(1, item.Time, item.Open) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
}

func httpGet(url string) (CryptoResponse, error) {
	var cryptoResponse CryptoResponse

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

func cryptoPipeline(fsym string, tsym string, limit int) {
	// Format url
	url := fmt.Sprintf(base_hour_hist_url + "?fsym=%s&tsym=%s&limit=%d", fsym, tsym, limit)
	var cryptoResponse CryptoResponse

	// Fetch data
	cryptoResponse, err := httpGet(url)
	if err != nil {
		log.Fatal("Failed http get request: ", err)
	}

	// Process response
	insertRecords(cryptoResponse.Data)
}


func main() {
	cryptoPipeline("BTC", "USD", 60)
}
