package models

type Record struct {
	Id            int       `json:"id"`
	MontorId      string    `json:"monitor_id"`
	RecordEpoch   string    `json:"record_epoch"`
	Value         float64   `json:"value"`
	Updated       string    `json:"updated"`
}

/*
	Inserts pricing data for given monitor id
 */
func InsertRecords(monitor_id int, data []HourHist) {
	// Construct prepared statement
	stmtIns, err := db.Prepare("INSERT INTO `records` (monitor_id, record_epoch, value, updated) VALUES(?,?,?,NOW())") // ? = placeholder
	if err != nil {
		panic(err.Error()) //TODO fix error handling
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Insert price records into db
	for i := range data {
		item := data[i]
		_, err = stmtIns.Exec(monitor_id, item.Time, item.Open) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) //TODO fix error handling
		}
	}
}
