package models

type Monitor struct {
    Id              int    `json:"id"`
    Name            string `json:"name"`
    Symbol          string `json:"symbol"`
    MonitorCategory int    `json:"monitor_category"`
    CreatedEpoch    int    `json:"created_epoch"`
    Enabled         bool   `json:"enabled"`
}

/*
    Fetches all active monitors from db
 */
func FetchMonitors() ([]*Monitor, error) {
    // Reading data
    rows, err := db.Query("SELECT * FROM monitors WHERE enabled = 1;")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    monitors := make([]*Monitor, 0)

    for rows.Next() {
        m := new(Monitor)
        err := rows.Scan(&m.Id, &m.Name, &m.Symbol, &m.MonitorCategory, &m.CreatedEpoch, &m.Enabled)
        if err != nil {
            return nil, err
        }
        monitors = append(monitors, m)
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }

    return monitors, nil
}