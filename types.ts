interface MinecraftChatMessage {
    message: string
    username: string
    uuid: string
    timestamp: string
    server: string
};

interface DiscordChatMessage {
    message: string
    username: string
    timestamp: string
    server: string
    channel_id: string
    guild_id: string
    guild_name: string
}


// type MinecraftChatAdvancement struct {
// 	Username    string `json:"username"`
// 	Advancement string `json:"advancement"`
// 	Time        int64  `json:"time"`
// 	MCServer    string `json:"mc_server"`
// 	UUID        string `json:"uuid,omitempty"`
// }

interface MinecraftAdvancementMessage { 
    advancement: string
    username: string
    uuid: string
    timestamp: string
    server: string
};

interface MinecraftPlayerJoinMessage {
    username: string
    uuid: string
    timestamp: string
    server: string
};

interface MinecraftPlayerLeaveMessage {
    username: string
    uuid: string
    timestamp: string
    server: string
};

interface MinecraftPlayerKillMessage {
    username: string
    uuid: string
    timestamp: string
    server: string
};

interface MinecraftPlayerDeathMessage {
    timestamp: string
    server: string

    victim: string
    death_message: string
    murderer?: string
    type: "pve" | "pvp"
    victimUUID: string
    murdererUUID?: string
};
