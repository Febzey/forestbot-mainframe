package database

import "time"

func (d *Database) INSERT_player_whois_description(username string, description string) error {

	now := time.Now().UnixMilli()

	INSERT_DESCRIPTION_QUERY := "INSERT INTO whois (username, description, timestamp) VALUES (?,?,?) ON DUPLICATE KEY UPDATE description = VALUES(description), timestamp = ?"
	_, err := d.Execute(INSERT_DESCRIPTION_QUERY, username, description, now, now)
	return err
}
