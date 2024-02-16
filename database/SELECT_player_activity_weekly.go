package database

import "time"

// +--------------+-------------+
// | COUNT(UUID)  | day_of_week |
// +--------------+-------------+
// | 15           | 1           |  -- Sunday
// | 25           | 2           |  -- Monday
// | 20           | 3           |  -- Tuesday
// | 18           | 4           |  -- Wednesday
// | 30           | 5           |  -- Thursday
// | 22           | 6           |  -- Friday
// | 17           | 7           |  -- Saturday
// +--------------+-------------+
// In this example:

// On Sunday (day_of_week = 1), 15 players logged in.
// On Monday (day_of_week = 2), 25 players logged in.

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
