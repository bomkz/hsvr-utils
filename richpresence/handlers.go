package richpresence

import (
	"bytes"
	"encoding/json"
	"fmt"

	"log"
	"strconv"
	"time"

	"github.com/bomkz/hsvr-utils/definitions"
	"github.com/google/uuid"
	"github.com/hugolgst/rich-go/client"
)

func DataTypeHandler(message bytes.Buffer, datatype string) {

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
		log.Print(err)
		go reconnectRPClient()
		reconnectingRP = true
	}
	if reconnectingRP {
		rpSuccess <- true
	}

}

func reconnectRPClient() {
	timer := time.NewTicker(10 * time.Second)
	for {

		if userOnline {
			select {
			case <-timer.C:
				if userOnline {
					handleUserOnline()
				} else {
					break
				}
			case <-rpSuccess:
				reconnectingRP = false
				break
			}
		} else {
			break
		}

	}
	timer.Stop()
}

func handleClientLogout() {
	client.Logout()

}

func rpHandler() {
	cooldown := time.NewTicker(16 * time.Second)
	for range cooldown.C {

		if userOnline {
			updateRichPresence()
		}

	}
}

func updateRichPresence() {

	if latestUserStats.Kills == 0 && latestUserStats.Deaths != 0 {
		latestUserStats.Ratio = float64(latestUserStats.Deaths) / -1
	} else if latestUserStats.Kills != 0 && latestUserStats.Deaths == 0 {
		latestUserStats.Ratio = float64(latestUserStats.Kills)
	} else {
		latestUserStats.Ratio = float64(latestUserStats.Kills) / float64(latestUserStats.Deaths)
	}

	kdr := strconv.Itoa(latestUserStats.Kills) + "K/" + strconv.Itoa(latestUserStats.Deaths) + "D/" + fmt.Sprint(latestUserStats.Ratio) + "R"
	state := "ELO: " + fmt.Sprint(latestUserStats.ELO) + " | " + kdr
	details := "VTOLVR 24/7RankedBVR"
	var aircraft string
	var smalltext string
	largetext := "Currently flying: "
	switch latestUserStats.CurrentVehicle {
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

	timejoin := time.Unix(latestUserStats.LastSpawnTimestamp.Unix(), 0)
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

func HandleInit() {
	go ConnectWS()

	go rpHandler()

	steamID32 := findCurrentUID()

	int64SteamID64 := convertID3ToID64(steamID32)

	steamID64 = strconv.Itoa(int(int64SteamID64))

	log.Println(steamID64)

}
func populateELO(message bytes.Buffer) {
	var newUserLookup definitions.UserLookupResultStruct

	err := json.Unmarshal(message.Bytes(), &newUserLookup)
	if err != nil {
		log.Fatal(err)
		return
	}

	latestUserStats.ELO = int(newUserLookup.Result.ELO)

}
func queryUser() {

	var newUserLookup definitions.LookupStruct

	newUUID, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
		return
	}
	newUserLookup.MessageType = "lookup"
	newUserLookup.Data.PID = newUUID.String()
	newUserLookup.Data.UID = steamID64
	newUserLookup.Data.Category = "user"

	newUserLookupByte, err := json.Marshal(newUserLookup)
	if err != nil {
		log.Println(err)
		return
	}

	err = localSocket.WriteString(string(newUserLookupByte))
	if err != nil {
		log.Println(err)
		return
	}
}

func checkIfUserIsOnline() bool {
	for _, y := range onlineUsers.Data {
		if y.UID == steamID64 {
			return true
		}
	}
	return false
}

func onUserLogin(message bytes.Buffer) {
	var UserLogin logStruct
	err := json.Unmarshal(message.Bytes(), &UserLogin)
	if err != nil {
		log.Println(err)
		return
	}

}

func onUserLogout(message bytes.Buffer) {
	var UserLogout logStruct
	err := json.Unmarshal(message.Bytes(), &UserLogout)
	if err != nil {
		log.Println(err)
		return
	}

}

func onOnlineUpdate(message bytes.Buffer) {
	var OnlineUpdate onlineStruct
	err := json.Unmarshal(message.Bytes(), &OnlineUpdate)
	if err != nil {
		log.Println(err)
		return
	}
	onlineUsers = OnlineUpdate

	isOnline := checkIfUserIsOnline()
	if isOnline && !userOnline {
		handleUserOnline()
	} else if !isOnline && userOnline {
		handleUserOffline()
	}

}

func onSpawn(message bytes.Buffer) {
	var Spawn spawnStruct
	err := json.Unmarshal(message.Bytes(), &Spawn)
	if err != nil {
		log.Println(err)
		return
	}

	for _, y := range Spawn.Data.User.Occupants {
		if y == steamID64 && !userOnline {
			handleUserOnline()
		}
		if y == steamID64 {
			latestUserStats.LastSpawnTimestamp = time.Now()
			latestUserStats.CurrentVehicle = Spawn.Data.User.Type
		}
	}

}

func onDeath(message bytes.Buffer) {
	var Death deathStruct
	err := json.Unmarshal(message.Bytes(), &Death)
	if err != nil {
		log.Println(err)
		return
	}

	for _, y := range Death.Data.Victim.Occupants {
		if y == steamID64 && !userOnline {
			handleUserOnline()
		}
		if y == steamID64 {
			latestUserStats.CurrentVehicle = "vtolvr"
			latestUserStats.Deaths += 1
			latestUserStats.SpawnedIn = false
		}
	}

}

func onKill(message bytes.Buffer) {
	var Kill killStruct
	err := json.Unmarshal(message.Bytes(), &Kill)
	if err != nil {
		log.Println(err)
		return
	}
	for _, y := range Kill.Data.Killer.Occupants {
		if y == steamID64 && !userOnline {
			handleUserOnline()
		}
		if y == steamID64 {
			latestUserStats.Kills += 1
		}

	}

}

func handleUserOnline() {
	userOnline = true
	handleClientLogin()
	queryUser()
}

func handleUserOffline() {
	latestUserStats = userStatsStruct{}
	handleClientLogout()
	userOnline = false
}
