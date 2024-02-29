package database

func (d *Database) UpdatePlayerPlaytime(uuid string, server string) error {
	_, err := d.Execute("UPDATE users SET playtime = playtime + 60000 WHERE uuid IN (?) AND mc_server = ?", uuid, server)
	if err != nil {
		return err
	}
	return nil
}
