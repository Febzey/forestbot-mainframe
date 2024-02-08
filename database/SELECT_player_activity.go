package database

import "github.com/febzey/ForestBot-Mainframe/types"

func (d *Database) GetAllPlayerActivity(server string, userOrUuid string, usingUuid bool) (interface{}, error) {

	SELECT_PLAYER_ACTIVITY := `
	SELECT
		UUID,
		Date,
		type
	FROM
		playerActivity
	WHERE
		username = ?
		AND mc_server = ?
		AND Date >= UNIX_TIMESTAMP(CURDATE() - INTERVAL 7 DAY) * 1000
		AND Date <= UNIX_TIMESTAMP(CURDATE()) * 1000
		AND (type = 'login' OR type = 'logout')
	ORDER BY
		Date;
	`

	if usingUuid {
		SELECT_PLAYER_ACTIVITY = `
		SELECT
			UUID,
			Date,
			type
		FROM
			playerActivity
		WHERE
			uuid = ?
			AND mc_server = ?
			AND Date >= UNIX_TIMESTAMP(CURDATE() - INTERVAL 7 DAY) * 1000
			AND Date <= UNIX_TIMESTAMP(CURDATE()) * 1000
			AND (type = 'login' OR type = 'logout')
		ORDER BY
			Date;
	`
	}
	rows, err := d.Query(SELECT_PLAYER_ACTIVITY, userOrUuid, server)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []types.PlayerActivity

	for rows.Next() {
		var activity types.PlayerActivity
		err := rows.Scan(
			&activity.UUID,
			&activity.Date,
			&activity.Type,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, activity)
	}
	return results, nil

}
