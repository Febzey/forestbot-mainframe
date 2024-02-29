package database

import (
	"fmt"

	"github.com/febzey/ForestBot-Mainframe/types"
)

func (d *Database) InsertPlayerDeathOrKill(args types.MinecraftPlayerDeathMessage) error {

	murderer := args.Murderer
	victim := args.Victim
	victim_uuid := args.VictimUUID
	murderer_uuid := args.MurdererUUID
	death_message := args.Death_message
	server := args.Mc_server
	time := args.Time

	fmt.Println("murderer:", murderer)
	fmt.Println("victim:", victim)
	fmt.Println("victim_uuid:", victim_uuid)
	fmt.Println("murderer_uuid:", murderer_uuid)
	fmt.Println("death_message:", death_message)
	fmt.Println("server:", server)
	fmt.Println("time:", time)

	//Updating the users death count, death message and death time.
	_, err := d.Execute("UPDATE users SET deaths = deaths + 1, lastdeathString = ?, lastdeathTime = ? WHERE username = ? AND mc_server = ?",
		death_message, time, victim, server,
	)
	if err != nil {
		return err
	}

	//If there is a murderer present.
	if murderer.Valid && murderer.String != "" {

		//Updating the murderers kill count.
		if _, err := d.Execute("UPDATE users SET kills = kills + 1 WHERE username = ? AND mc_server = ?", murderer.String, server); err != nil {
			return err
		}

		//Inserting the death into deaths table. with murderer and victim as type PVP
		if _, err := d.Execute("INSERT into deaths (victim, death_message, murderer, time, type, mc_server, victimUUID, murdererUUID) VALUES (?,?,?,?,?,?,?,?)",
			victim, death_message, murderer.String, time, "pvp", server, victim_uuid, murderer_uuid.String); err != nil {
			return err
		}

		//No murderer
	} else {

		//No murderer was found so save to deaths table as PVE death
		if _, err := d.Execute("INSERT into deaths (victim, death_message, time, type, mc_server, victimUUID) VALUES (?, ?, ?, ?, ?, ?)",
			victim, death_message, time, "pve", server, victim_uuid); err != nil {
			return err
		}

	}

	return nil
}
