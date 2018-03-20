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

type Record struct {
	Id            int       `json:"id"`
	MontorId      string    `json:"monitor_id"`
	RecordEpoch   string    `json:"record_epoch"`
	Value         float64   `json:"value"`
	Updated       string    `json:"updated"`
}

type Monitor struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	MonitorCategory int    `json:"monitor_category"`
	CreatedEpoch    int    `json:"created_epoch"`
	Enabled         bool   `json:"enabled"`
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

/*
	Inserts pricing data for given monitor id
 */
func insertRecords(monitor_id int, data []HourHist) {
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
		_, err = stmtIns.Exec(monitor_id, item.Time, item.Open) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
}

/*
	Fetches all active monitors from db
 */
func fetchMonitors() []Monitor {
	db, err := sql.Open("mysql", "root:my-secret-pw@tcp(localhost:13306)/clurichaun")
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Reading data
	rows, err := db.Query("SELECT * FROM monitors WHERE enabled = 1;")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	var id, monitorCategory, createdEpoch int
	var name, symbol string
	var enabled bool
	var monitors []Monitor

	for rows.Next() {
		err := rows.Scan(&id, &name, &symbol, &monitorCategory, &createdEpoch, &enabled)
		if err != nil {
			panic(err.Error())
		}
		monitors = append(monitors, Monitor{Id: id, Name: name, Symbol: symbol, MonitorCategory: monitorCategory, CreatedEpoch: createdEpoch, Enabled: enabled})
	}

	return monitors
}

/*
	Fetches crypto price data given url
 */
func cryptoGet(url string) (CryptoResponse, error) {

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

/*
	Runs crypto pipeline backfill process for each monitor for specified time range
 */
func cryptoPipeline(monitors []Monitor, hoursBack int) {
	// Iterate over monitors
	for _, monitor := range monitors {
		// Format url
		url := fmt.Sprintf(base_hour_hist_url + "?fsym=%s&tsym=%s&limit=%d", monitor.Symbol, "USD", hoursBack)
		var cryptoResponse CryptoResponse

		// Fetch data
		cryptoResponse, err := cryptoGet(url)
		if err != nil {
			log.Fatal("Failed http get request: ", err)
		}

		// Process response
		insertRecords(monitor.Id, cryptoResponse.Data)
	}
}

func main() {
	monitors := fetchMonitors()
	cryptoPipeline(monitors, 60)
}
