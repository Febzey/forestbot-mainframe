package database

type ServerStatsProps struct {
	TotalLogins   int
	UniquePlayers int
	UniqueLogins  int

	UserWithMostLogins struct {
		Username   string
		LoginCount int
	}

	TotalUsersSaved   int
	TotalAdvancements int
	TotalDeaths       int
}

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
	AND date >= UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 7 DAY)) * 1000
	GROUP BY username
	ORDER BY login_count DESC
	LIMIT 1;
	`

	SELECT_TOTAL_USERS_SAVED = `
	SELECT COUNT(*) AS user_count
	FROM users
	WHERE mc_server = ?;
	`

	SELECT_TOTAL_ADVANCEMENTS = `
	SELECT COUNT(*) AS advancement_count
	FROM advancements
	WHERE mc_server = ?;
	`

	SELECT_TOTAL_DEATHS = `
	SELECT COUNT(*) AS death_count
	FROM deaths
	WHERE mc_server = ?;
	`
)

// some weekly and overall stats for the server
func (d *Database) SELECT_server_stats(server string) (stats ServerStatsProps, err error) {

	var totalLogins int
	err = d.Pool.QueryRow(SELECT_TOTAL_LOGINS, server).Scan(&totalLogins)
	if err != nil {
		return stats, err
	}

	var totalUniqueLogins int
	err = d.Pool.QueryRow(SELECT_TOTAL_UNIQUE_LOGINS, server).Scan(&totalUniqueLogins)
	if err != nil {
		return stats, err
	}

	var totalNewUsers int
	err = d.Pool.QueryRow(SELECT_TOTAL_NEW_USERS_COUNT, server).Scan(&totalNewUsers)
	if err != nil {
		return stats, err
	}

	err = d.Pool.QueryRow(SELECT_USER_WITH_MOST_LOGINS, server).Scan(&stats.UserWithMostLogins.Username, &stats.UserWithMostLogins.LoginCount)
	if err != nil {
		return stats, err
	}

	err = d.Pool.QueryRow(SELECT_TOTAL_USERS_SAVED, server).Scan(&stats.TotalUsersSaved)
	if err != nil {
		return stats, err
	}

	err = d.Pool.QueryRow(SELECT_TOTAL_ADVANCEMENTS, server).Scan(&stats.TotalAdvancements)
	if err != nil {
		return stats, err
	}

	err = d.Pool.QueryRow(SELECT_TOTAL_DEATHS, server).Scan(&stats.TotalDeaths)
	if err != nil {
		return stats, err
	}

	stats = ServerStatsProps{
		TotalLogins:        totalLogins,
		UniquePlayers:      totalUniqueLogins,
		UniqueLogins:       totalNewUsers,
		UserWithMostLogins: stats.UserWithMostLogins,
		TotalUsersSaved:    stats.TotalUsersSaved,
		TotalAdvancements:  stats.TotalAdvancements,
		TotalDeaths:        stats.TotalDeaths,
	}

	return stats, err

}

var (
	// top 5 pvpers
	SELECT_TOP_5_KILLERS = `
    SELECT murderer AS player_name, murdererUUID AS player_uuid,
    COUNT(*) AS kill_count
    FROM deaths
    WHERE type = 'pvp'
    AND mc_server = ?
    AND time >= UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 7 DAY)) * 1000
    GROUP BY murderer
    ORDER BY kill_count DESC
    LIMIT 5;
    `

	// top 5 pve deaths
	SELECT_TOP_5_PVE_DEATHS = `
    SELECT victim AS player_name, victimUUID AS player_uuid,
    COUNT(*) AS death_count
    FROM deaths
    WHERE type = 'pve'
    AND time >= UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 7 DAY)) * 1000
    AND mc_server = ?
    GROUP BY victim
    ORDER BY death_count DESC
    LIMIT 5;
    `

	// top 5 pvp deaths
	SELECT_TOP_5_PVP_DEATHS = `
    SELECT victim AS player_name, victimUUID AS player_uuid,
    COUNT(*) AS pvp_death_count
    FROM deaths
    WHERE type = 'pvp'
    AND time >= UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 7 DAY)) * 1000
    AND mc_server = ?
    GROUP BY victim
    ORDER BY pvp_death_count DESC
    LIMIT 5;
    `

	// top 5 users who got most advancements
	SELECT_TOP_5_ADVANCEMENTS = `
    SELECT username AS player_name, uuid AS player_uuid,
    COUNT(*) AS advancement_count
    FROM advancements
    WHERE time >= UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 7 DAY)) * 1000
    AND mc_server = ?
    GROUP BY username
    ORDER BY advancement_count DESC
    LIMIT 5;
    `

	// top 5 players with most logins
	SELECT_TOP_5_LOGINS = `
    SELECT username, uuid AS player_uuid, COUNT(*) AS login_count
    FROM playerActivity
    WHERE type = 'login'
    AND date >= UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 7 DAY)) * 1000
    AND mc_server = ?
    GROUP BY username
    ORDER BY login_count DESC
    LIMIT 5;
    `
)

type top5Leaderboards struct {
	Top5Killers []struct {
		PlayerName string
		KillCount  int
		Uuid       string
	}

	Top5PVEDeaths []struct {
		PlayerName string
		DeathCount int
		Uuid       string
	}

	Top5PVPDeaths []struct {
		PlayerName    string
		PVPDeathCount int
		Uuid          string
	}

	Top5Advancements []struct {
		PlayerName       string
		AdvancementCount int
		Uuid             string
	}

	Top5Logins []struct {
		PlayerName string
		LoginCount int
		Uuid       string
	}
}

// GetTop5Leaderboard is a function that returns the top 5 leaderboards for various things from the past 7 days exactly.
func (d *Database) SELECT_top_5_player_stats(server string) (stats top5Leaderboards, err error) {
	var top5 top5Leaderboards

	// Top 5 PVP Kills
	rows, err := d.Pool.Query(SELECT_TOP_5_KILLERS, server)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var top5Killer struct {
			PlayerName string
			KillCount  int
			Uuid       string
		}

		err := rows.Scan(&top5Killer.PlayerName, &top5Killer.Uuid, &top5Killer.KillCount)
		if err != nil {
			return stats, err
		}

		top5.Top5Killers = append(top5.Top5Killers, top5Killer)

	}
	if err := rows.Err(); err != nil {
		return stats, err
	}

	// Top 5 PVE Deaths
	rows, err = d.Pool.Query(SELECT_TOP_5_PVE_DEATHS, server)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var top5PVEDeath struct {
			PlayerName string
			DeathCount int
			Uuid       string
		}

		err := rows.Scan(&top5PVEDeath.PlayerName, &top5PVEDeath.Uuid, &top5PVEDeath.DeathCount)
		if err != nil {
			return stats, err
		}

		top5.Top5PVEDeaths = append(top5.Top5PVEDeaths, top5PVEDeath)

	}
	if err := rows.Err(); err != nil {
		return stats, err
	}

	// Top 5 PVP Deaths
	rows, err = d.Pool.Query(SELECT_TOP_5_PVP_DEATHS, server)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var top5PVPDeath struct {
			PlayerName    string
			PVPDeathCount int
			Uuid          string
		}

		err := rows.Scan(&top5PVPDeath.PlayerName, &top5PVPDeath.Uuid, &top5PVPDeath.PVPDeathCount)
		if err != nil {
			return stats, err
		}

		top5.Top5PVPDeaths = append(top5.Top5PVPDeaths, top5PVPDeath)
	}
	if err := rows.Err(); err != nil {
		return stats, err
	}

	// Top 5 Advancements
	rows, err = d.Pool.Query(SELECT_TOP_5_ADVANCEMENTS, server)
	if err != nil {
		return stats, err
	}

	defer rows.Close()

	for rows.Next() {
		var top5Advancement struct {
			PlayerName       string
			AdvancementCount int
			Uuid             string
		}

		err := rows.Scan(&top5Advancement.PlayerName, &top5Advancement.Uuid, &top5Advancement.AdvancementCount)
		if err != nil {
			return stats, err
		}

		top5.Top5Advancements = append(top5.Top5Advancements, top5Advancement)
	}
	if err := rows.Err(); err != nil {
		return stats, err
	}

	// Top 5 Logins
	rows, err = d.Pool.Query(SELECT_TOP_5_LOGINS, server)
	if err != nil {
		return stats, err
	}

	defer rows.Close()

	for rows.Next() {
		var top5Login struct {
			PlayerName string
			LoginCount int
			Uuid       string
		}

		err := rows.Scan(&top5Login.PlayerName, &top5Login.Uuid, &top5Login.LoginCount)
		if err != nil {
			return stats, err
		}

		top5.Top5Logins = append(top5.Top5Logins, top5Login)
	}

	if err := rows.Err(); err != nil {
		return stats, err
	}

	return top5, err
}
