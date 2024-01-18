package types

type PlayerList struct {
	Name    string
	Ping    int
	HeadURL string
}

type User struct {
	Username        string
	Kills           int64
	Deaths          int64
	Joindate        string
	LastSeen        string
	UUID            string
	Playtime        int64
	Joins           int64
	Leaves          int64
	LastDeathTime   int64
	LastDeathString string
	MCServer        string
}
