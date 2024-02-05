package database

// PlayerActivityHourlyResult represents a single result entry for the hourly player activity.
type PlayerActivityHourlyResult struct {
	Weekday  int
	Activity []HourlyActivity
}

// HourlyActivity represents the hourly activity data.
type HourlyActivity struct {
	Hour   int
	Logins int
}

func (c *Database) PlayerActivityHourlyResults(server string) ([]PlayerActivityHourlyResult, error) {

	var SELECT_HOURLY_PLAYER_ACTIVITY = `
	SELECT
		COUNT(DISTINCT UUID) AS user_count,
		DAYOFWEEK(FROM_UNIXTIME(Date / 1000)) AS day_of_week,
		HOUR(FROM_UNIXTIME(Date / 1000)) AS hour_of_day
	FROM playerActivity
	WHERE mc_server = ?
		AND Date >= UNIX_TIMESTAMP(CURDATE() - INTERVAL 7 DAY) * 1000
		AND Date <= UNIX_TIMESTAMP(CURDATE()) * 1000
		AND HOUR(FROM_UNIXTIME(Date / 1000)) BETWEEN 0 AND 23
		AND type = "login"
	GROUP BY day_of_week, hour_of_day
	`

	rows, err := c.Query(SELECT_HOURLY_PLAYER_ACTIVITY, server)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]PlayerActivityHourlyResult, 0)

	for rows.Next() {
		var userCount, dayOfWeek, hourOfDay int
		err := rows.Scan(&userCount, &dayOfWeek, &hourOfDay)
		if err != nil {
			return nil, err
		}

		// Find existing entry in results
		var foundEntry *PlayerActivityHourlyResult
		for i, entry := range results {
			if entry.Weekday == dayOfWeek {
				foundEntry = &results[i]
				break
			}
		}

		// If not found, create a new entry
		if foundEntry == nil {
			foundEntry = &PlayerActivityHourlyResult{
				Weekday:  dayOfWeek,
				Activity: make([]HourlyActivity, 0),
			}
			results = append(results, *foundEntry)
		}

		// Append the hourly activity
		foundEntry.Activity = append(foundEntry.Activity, HourlyActivity{
			Hour:   hourOfDay,
			Logins: userCount,
		})
	}

	return results, nil
	// Print the JSON for debugging
}
