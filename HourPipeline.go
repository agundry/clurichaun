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

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO `records` (monitor_id, record_epoch, value, updated) VALUES(?,?,?,NOW())") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Insert square numbers for 0-24 in the database
	for i := range data {
		item := data[i]
		_, err = stmtIns.Exec(1, item.Time, item.Open) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	// 	fmt.Println(cryptoResponse.Data[i].Time + '\n')  // This prints '0', two times 
	// }
}

func main() {
	// phone := "14158586273"
	// QueryEscape escapes the phone string so
	// it can be safely placed inside a URL query
	// safePhone := url.QueryEscape(phone)

	url := fmt.Sprintf(base_hour_hist_url + "?fsym=%s&tsym=%s&limit=%d", "BTC", "USD", 60)
	fmt.Println(url)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Use json.Decode for reading streams of JSON data
	var cryptoResponse CryptoResponse
	if err := json.NewDecoder(resp.Body).Decode(&cryptoResponse); err != nil {
		log.Println(err)
	}

	fmt.Println(cryptoResponse.Data[0].Time)
	insertRecords(cryptoResponse.Data)

	// for i := range cryptoResponse.Data {
	// 	fmt.Println(cryptoResponse.Data[i].Time + '\n')  // This prints '0', two times 
	// }

	// fmt.Println(cryptoResponse)
	// fmt.Println("Data = ", cryptoResponse.Data)
	// fmt.Println("Country   = ", record.CountryName)
	// fmt.Println("Location  = ", record.Location)
	// fmt.Println("Carrier   = ", record.Carrier)
	// fmt.Println("LineType  = ", record.LineType)

}
