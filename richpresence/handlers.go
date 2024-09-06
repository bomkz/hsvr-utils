package richpresence

import (
	"bytes"
	"encoding/json"
	"fmt"

	"log"
	"strconv"
	"time"

	"github.com/bomkz/vtolvr-utils/definitions"
	"github.com/google/uuid"
	"github.com/hugolgst/rich-go/client"
)

func DataTypeHandler(message bytes.Buffer, datatype string) {

	log.Println(datatype+": ", message.String())
	switch datatype {
	case "user_logout":
		onUserLogout(message)
	case "user_login":
		onUserLogin(message)
	case "online":
		onOnlineUpdate(message)
	case "spawn":
		onSpawn(message)
	case "death":
		onDeath(message)
	case "kill":
		onKill(message)
	case "response":
		populateELO(message)
	}
}

func handleClientLogin() {
	err := client.Login("1220960048704913448")
	if err != nil {
		log.Fatal(err)
	}
}

func handleClientLogout() {
	client.Logout()

}

func RichPresenceHandler() {
	cooldown := time.NewTicker(16 * time.Second)
	for range cooldown.C {

		if definitions.UserOnline {
			updateRichPresence()
		}

	}
}

func updateRichPresence() {
	definitions.LatestUserStats.Ratio = float64(definitions.LatestUserStats.Kills) / float64(definitions.LatestUserStats.Deaths)

	kdr := strconv.Itoa(definitions.LatestUserStats.Kills) + "K/" + strconv.Itoa(definitions.LatestUserStats.Deaths) + "D/" + fmt.Sprint(definitions.LatestUserStats.Ratio) + "R"
	state := "ELO: " + fmt.Sprint(definitions.LatestUserStats.ELO) + " | " + kdr
	details := "VTOLVR 24/7RankedBVR"
	var aircraft string
	var smalltext string
	largetext := "Currently flying: "
	switch definitions.LatestUserStats.CurrentVehicle {
	case "vtolvr":
		largetext += "In Lobby"
		aircraft = "vtolvr"
		smalltext = "It's the hit game CTOL XR by BahamutoC!!"
	case "Vehicles/EF-24":
		largetext += "EF-24G Mischief"
		aircraft = "ef24g"
		smalltext = "Split the throttles soulja boy!"
	case "Vehicles/SEVTF":
		largetext += "F-45A Ghost"
		aircraft = "f45a"
		smalltext = "The mind of an f45 main cannot comprehend a 26 chaffing."
	case "Vehicles/FA-26B":
		largetext += "F/A-26B Wasp"
		aircraft = "fa26b"
		smalltext = "Carrying literally the entire weight of a Shipping Container in bombs."
	case "Vehicles/T-55":
		largetext += "T-55 Tyro"
		aircraft = "t55"
		smalltext = "Courtesy of dubyaaa"
	}

	timejoin := time.Unix(definitions.LatestUserStats.LastSpawnTimestamp.Unix(), 0)
	err := client.SetActivity(client.Activity{
		State:      state,
		Details:    details,
		LargeImage: aircraft,
		LargeText:  largetext,
		SmallImage: "vtolvr",
		SmallText:  smalltext,

		Timestamps: &client.Timestamps{
			Start: &timejoin,
		},
	})
	if err != nil {
		log.Println(err)
	}
}

func populateELO(message bytes.Buffer) {
	var newUserLookup definitions.UserLookupResultStruct

	err := json.Unmarshal(message.Bytes(), &newUserLookup)
	if err != nil {
		log.Fatal(err)
		return
	}

	definitions.LatestUserStats.ELO = int(newUserLookup.Result.ELO)

}
func queryUser() {

	var newUserLookup definitions.UserLookup

	newUUID, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
		return
	}
	newUserLookup.MessageType = "lookup"
	newUserLookup.Data.PID = newUUID.String()
	newUserLookup.Data.UID = definitions.SteamID64
	newUserLookup.Data.Category = "user"

	newUserLookupByte, err := json.Marshal(newUserLookup)
	if err != nil {
		log.Println(err)
		return
	}

	err = definitions.Socket.WriteString(string(newUserLookupByte))
	if err != nil {
		log.Println(err)
		return
	}
}

func checkIfUserIsOnline() bool {
	for _, y := range definitions.OnlineUsers.Data {
		if y.UID == definitions.SteamID64 {
			return true
		}
	}
	return false
}

func onUserLogin(message bytes.Buffer) {
	var UserLogin definitions.LogStruct
	err := json.Unmarshal(message.Bytes(), &UserLogin)
	if err != nil {
		log.Println(err)
		return
	}

}

func onUserLogout(message bytes.Buffer) {
	var UserLogout definitions.LogStruct
	err := json.Unmarshal(message.Bytes(), &UserLogout)
	if err != nil {
		log.Println(err)
		return
	}

}

func onOnlineUpdate(message bytes.Buffer) {
	var OnlineUpdate definitions.OnlineStruct
	err := json.Unmarshal(message.Bytes(), &OnlineUpdate)
	if err != nil {
		log.Println(err)
		return
	}
	definitions.OnlineUsers = OnlineUpdate

	isOnline := checkIfUserIsOnline()
	if isOnline && !definitions.UserOnline {
		handleUserOnline()
	} else if !isOnline && definitions.UserOnline {
		handleUserOffline()
	}

}

func onSpawn(message bytes.Buffer) {
	var Spawn definitions.SpawnStruct
	err := json.Unmarshal(message.Bytes(), &Spawn)
	if err != nil {
		log.Println(err)
		return
	}

	for _, y := range Spawn.Data.User.Occupants {
		if y == definitions.SteamID64 && !definitions.UserOnline {
			handleUserOnline()
		}
		if y == definitions.SteamID64 {
			definitions.LatestUserStats.LastSpawnTimestamp = time.Now()
			definitions.LatestUserStats.CurrentVehicle = Spawn.Data.User.Type
		}
	}

}

func onDeath(message bytes.Buffer) {
	var Death definitions.DeathStruct
	err := json.Unmarshal(message.Bytes(), &Death)
	if err != nil {
		log.Println(err)
		return
	}

	for _, y := range Death.Data.Victim.Occupants {
		if y == definitions.SteamID64 && !definitions.UserOnline {
			handleUserOnline()
		}
		if y == definitions.SteamID64 {
			definitions.LatestUserStats.CurrentVehicle = "vtolvr"
			definitions.LatestUserStats.Deaths += 1
			definitions.LatestUserStats.SpawnedIn = false
		}
	}

}

func onKill(message bytes.Buffer) {
	var Kill definitions.KillStruct
	err := json.Unmarshal(message.Bytes(), &Kill)
	if err != nil {
		log.Println(err)
		return
	}
	for _, y := range Kill.Data.Killer.Occupants {
		if y == definitions.SteamID64 && !definitions.UserOnline {
			handleUserOnline()
		}
		if y == definitions.SteamID64 {
			definitions.LatestUserStats.Kills += 1
		}

	}

}

func handleUserOnline() {
	definitions.UserOnline = true
	handleClientLogin()
	queryUser()
	handleClientLogin()
}

func handleUserOffline() {
	definitions.LatestUserStats = definitions.UserStatsStruct{}
	handleClientLogout()
	definitions.UserOnline = false
}
