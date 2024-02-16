package database

func (d *Database) UniqueServers() ([]string, error) {
	var uniqueServers []string

	// Assuming your table is named 'your_table_name' and the column is 'mc_server'
	query := "SELECT DISTINCT mc_server FROM users WHERE mc_server IS NOT NULL AND LENGTH(mc_server) > 0"

	rows, err := d.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var server string
		if err := rows.Scan(&server); err != nil {
			return nil, err
		}
		uniqueServers = append(uniqueServers, server)
	}

	return uniqueServers, nil
}
