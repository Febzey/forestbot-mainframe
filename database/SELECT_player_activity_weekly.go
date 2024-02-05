package database

import "time"

// Total number of logins for each day of the week
func (db *Database) PlayerActivityWeekResults(mc_server string) (map[string]int, error) {

	rows, err := db.Query(`
	SELECT COUNT(UUID), DAYOFWEEK(FROM_UNIXTIME(Date / 1000)) AS day_of_week
	FROM playerActivity
	WHERE mc_server = ?
	  AND Date >= UNIX_TIMESTAMP(CURDATE() - INTERVAL 7 DAY) * 1000
	  AND Date <= UNIX_TIMESTAMP(CURDATE()) * 1000
	  AND type = "login"
	GROUP BY day_of_week
	`, mc_server)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make(map[string]int)

	for rows.Next() {
		var count int
		var day int
		err := rows.Scan(&count, &day)
		if err != nil {
			return nil, err
		}

		weekdayStr := time.Weekday(day).String()
		results[weekdayStr] = count

	}

	return results, nil

}
