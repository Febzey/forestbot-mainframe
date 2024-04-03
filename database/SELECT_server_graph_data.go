package database

import "fmt"

type PlayerActivityHourlyResult struct {
	Weekday  int
	Activity []HourlyActivity
}

type ServerStats struct {
	PlayerActivityHourlyResult []PlayerActivityHourlyResult

	TotalLogins        int
	UniquePlayers      int
	UniqueLogins       int
	UserWithMostLogins struct {
		Username   string
		LoginCount int
	}
}

type HourlyActivity struct {
	Hour   int
	Logins int
}

// type ServerStatsProps struct {

// }

var (
	SELECT_TOTAL_LOGINS = `
	SELECT COUNT(*) AS unique_logins_count
	FROM playerActivity
	WHERE type = 'login'
	AND mc_server = ?
	AND date >= UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 7 DAY)) * 1000;
	`

	SELECT_TOTAL_UNIQUE_LOGINS = `
	SELECT COUNT(DISTINCT UUID) AS unique_logins_count
	FROM playerActivity
	WHERE type = 'login'
	AND mc_server = ?
	AND date >= UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 7 DAY)) * 1000;
	`

	SELECT_TOTAL_NEW_USERS_COUNT = `
	SELECT COUNT(*) AS new_players_count
	FROM (
		SELECT UUID
		FROM playerActivity
		WHERE type = 'login'
		AND mc_server = ?
		GROUP BY UUID
		HAVING MIN(date) >= UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 7 DAY)) * 1000
	) AS new_players;
	`

	//! TODO lets remove this query. since we already send over the top 5 players with most logins for the past 7 days
	//! remove once updated in frontend
	SELECT_USER_WITH_MOST_LOGINS = `
	SELECT username, COUNT(*) AS login_count
	FROM playerActivity
	WHERE type = 'login'
	AND mc_server = ?
	AND date >= UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 10 DAY)) * 1000
	GROUP BY username
	ORDER BY login_count DESC
	LIMIT 1;
	`
)

// +------------+-------------+--------------+
// | user_count | day_of_week | hour_of_day  |
// +------------+-------------+--------------+
// | 5          | 1           | 12           |  -- Sunday, 12:00 AM
// | 8          | 1           | 13           |  -- Sunday, 1:00 PM
// | 10         | 2           | 14           |  -- Monday, 2:00 PM
// | ...        | ...         | ...          |
// +------------+-------------+--------------+

// This table indicates, for example, that on Sunday at 12:00 AM, there were 5 unique players who logged in, and on Sunday at 1:00 PM, there were 8 unique players who logged in. The counts are broken down by the day of the week and hour of the day.

func (d *Database) ServerActivityHourlyResults(server string) (stats ServerStats, err error) {

	var SELECT_HOURLY_PLAYER_ACTIVITY = `
	SELECT
		COUNT(DISTINCT UUID) AS user_count,
		DAYOFWEEK(FROM_UNIXTIME(Date / 1000)) AS day_of_week,
		HOUR(FROM_UNIXTIME(Date / 1000)) AS hour_of_day
	FROM playerActivity
	WHERE mc_server = ?
		AND Date >= UNIX_TIMESTAMP(CURDATE() - INTERVAL 10 DAY) * 1000
		AND Date <= UNIX_TIMESTAMP(CURDATE()) * 1000
		AND HOUR(FROM_UNIXTIME(Date / 1000)) BETWEEN 0 AND 23
		AND type = "login"
	GROUP BY day_of_week, hour_of_day
	`

	rows, err := d.Query(SELECT_HOURLY_PLAYER_ACTIVITY, server)
	if err != nil {
		fmt.Println(err, " heree")
		return stats, err
	}

	defer rows.Close()

	results := make([]PlayerActivityHourlyResult, 0)

	for rows.Next() {
		var userCount, dayOfWeek, hourOfDay int
		err := rows.Scan(&userCount, &dayOfWeek, &hourOfDay)
		if err != nil {
			fmt.Println(err)
			return stats, err
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

	var totalLogins int
	err = d.Pool.QueryRow(SELECT_TOTAL_LOGINS, server).Scan(&totalLogins)
	if err != nil {
		fmt.Println(err, " Error in SELECT_TOTAL_LOGINS")
		return stats, err
	}

	var totalUniqueLogins int
	err = d.Pool.QueryRow(SELECT_TOTAL_UNIQUE_LOGINS, server).Scan(&totalUniqueLogins)
	if err != nil {
		fmt.Println(err, " Error in SELECT_TOTAL_UNIQUE_LOGINS")
		return stats, err
	}

	var totalNewUsers int
	err = d.Pool.QueryRow(SELECT_TOTAL_NEW_USERS_COUNT, server).Scan(&totalNewUsers)
	if err != nil {
		fmt.Println(err, " Error in SELECT_TOTAL_NEW_USERS_COUNT")
		return stats, err
	}

	err = d.Pool.QueryRow(SELECT_USER_WITH_MOST_LOGINS, server).Scan(&stats.UserWithMostLogins.Username, &stats.UserWithMostLogins.LoginCount)
	if err != nil {
		fmt.Println(err, " Error in SELECT_USER_WITH_MOST_LOGINS")
		return stats, err
	}

	stats = ServerStats{
		PlayerActivityHourlyResult: results,
		TotalLogins:                totalLogins,
		UniquePlayers:              totalUniqueLogins,
		UniqueLogins:               totalNewUsers,
		UserWithMostLogins:         stats.UserWithMostLogins,
	}

	return stats, nil
	// Print the JSON for debugging
}
