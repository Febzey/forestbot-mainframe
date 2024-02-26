package types

import (
	"database/sql"
	"image"
	"sync"
)

type ImageCache struct {
	//The map of player heads to their image in the tablist.
	HeadImages map[string]image.Image

	//The mutex for the head images map.
	Mu sync.RWMutex
}

type PlayerActivity struct {
	ID        int    //id
	UUID      string //uuid
	Username  string //username
	Date      int64  // unix timestamp in milliseconds
	Type      string // join | leave
	Mc_server string // mc_server
}

type User struct {
	Username        string
	Kills           int64
	Deaths          int64
	Joindate        string
	LastSeen        sql.NullString
	UUID            sql.NullString
	Playtime        int64
	Joins           int64
	Leaves          int64
	LastDeathTime   int64
	LastDeathString sql.NullString
	MCServer        string
}

type Player struct {
	Username string `json:"username"`
	Uuid     string `json:"uuid"`
	Latency  int    `json:"latency"`
	Server   string `json:"server"`
	Head_url string `json:"head_url"`
}

type DiscordMessage struct {
	Message    string `json:"message"`
	Username   string `json:"username"`
	Timestamp  string `json:"timestamp"`
	Server     string `json:"server"`
	Channel_ID string `json:"channel_id"`
	Guild_ID   string `json:"guild_id"`
	Guild_Name string `json:"guild_name"`
}

type MinecraftChatMessage struct {
	Name      string         `json:"name"`
	Message   string         `json:"message"`
	Date      sql.NullString `json:"date"`
	Mc_server string         `json:"mc_server"`
	Uuid      string         `json:"uuid"`
}

type MinecraftAdvancementMessage struct {
	Username    string `json:"username"`
	Advancement string `json:"advancement"`
	Time        int64  `json:"time"`
	Mc_server   string `json:"mc_server"`
	Id          int    `json:"id,omitempty"`
	Uuid        string `json:"uuid"`
}

type MinecraftPlayerJoinMessage struct {
	Username  string `json:"username"`
	Uuid      string `json:"uuid"`
	Timestamp string `json:"timestamp"`
	Server    string `json:"server"`
	Latency   int    `json:"latency"`
}

type MinecraftPlayerLeaveMessage struct {
	Username  string `json:"username"`
	Uuid      string `json:"uuid"`
	Timestamp string `json:"timestamp"`
	Server    string `json:"server"`
}

type MinecraftPlayerKillMessage struct {
	Username  string `json:"username"`
	Uuid      string `json:"uuid"`
	Timestamp string `json:"timestamp"`
	Server    string `json:"server"`
}

type MinecraftPlayerDeathMessage struct {
	Victim        string          `json:"victim"`
	Death_message string          `json:"death_message"`
	Murderer      *sql.NullString `json:"murderer,omitempty"`
	Time          int64           `json:"time"`
	Type          string          `json:"type"`
	Mc_server     string          `json:"mc_server"`
	Id            int
	VictimUUID    string          `json:"victimUUID"`
	MurdererUUID  *sql.NullString `json:"murdererUUID"`
}
