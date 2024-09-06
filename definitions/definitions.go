package definitions

import (
	_ "embed"
	"time"

	"fyne.io/fyne/v2"
	"github.com/google/uuid"
	"github.com/lxzan/gws"
)

//go:embed .\..\aircraft.ico
var Icon []byte

type PongStruct struct {
	PID         uuid.UUID `json:"pid"`
	MessageType string    `json:"type"`
}
type PayloadTypeStruct struct {
	MessageType string `json:"type"`
}

type SubscribeStruct struct {
	MessageType string   `json:"type"`
	Data        []string `json:"data"`
}

type OnlineStruct struct {
	PID         uuid.UUID          `json:"pid"`
	MessageType string             `json:"type"`
	Data        []OnlineDataStruct `json:"data"`
}

type OnlineDataStruct struct {
	Name string `json:"name"`
	UID  string `json:"id"`
	Team string `json:"team"`
}

type LogStruct struct {
	PID         uuid.UUID     `json:"pid"`
	MessageType string        `json:"type"`
	Data        LogDataStruct `json:"data"`
}

type LogDataStruct struct {
	UserID    string `json:"userId"`
	PilotName string `json:"pilotName,omitempty"`
}

type SpawnStruct struct {
	PID         uuid.UUID       `json:"pid"`
	MessageType string          `json:"type"`
	Data        SpawnDataStruct `json:"data"`
}

type SpawnDataStruct struct {
	User       UserStruct       `json:"user"`
	ServerInfo ServerInfoStruct `json:"serverInfo"`
}

type ServerInfoStruct struct {
	OnlineUsers []string          `json:"onlineUsers"`
	MissionID   string            `json:"missionId"`
	Environment EnvironmentStruct `json:"environment"`
}

type EnvironmentStruct struct {
	TimeOfDay float64    `json:"tod"`
	Weather   int        `json:"weather"`
	Wind      WindStruct `json:"wind"`
}

type WindStruct struct {
	Heading   float64 `json:"heading"`
	Magnitude float64 `json:"mag"`
	Variance  float64 `json:"var"`
	Gusts     float64 `json:"gust"`
}

type UserStruct struct {
	OwnerID   string   `json:"ownerId"`
	Occupants []string `json:"occupants"`
	Position  XYZ      `json:"position"`
	Velocity  XYZ      `json:"velocity"`
	Team      string   `json:"team"`
	Type      string   `json:"type"`
}

type XYZ struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type KillStruct struct {
	PID         uuid.UUID      `json:"pid"`
	MessageType string         `json:"type"`
	Data        KillDataStruct `json:"data"`
}
type KillDataStruct struct {
	Victim                  UserStruct       `json:"victim"`
	Killer                  UserStruct       `json:"killer"`
	ServerInfo              ServerInfoStruct `json:"serverInfo"`
	Weapon                  string           `json:"weapon"`
	WeaponUUID              string           `json:"WeaponUuid"`
	PreviousDamagedByUserId string           `json:"previousDamagedByUserId"`
	PreviousDamagedByWeapon string           `json:"previousDamagedByWeapon"`
}

type DeathStruct struct {
	PID         string          `json:"pid"`
	MessageType string          `json:"type"`
	Data        DeathDataStruct `json:"data"`
}

type DeathDataStruct struct {
	Victim     UserStruct       `json:"victim"`
	ServerInfo ServerInfoStruct `json:"serverInfo"`
}

type UserLookup struct {
	MessageType string         `json:"type"`
	Data        UserLookupData `json:"data"`
}

type UserLookupData struct {
	UID      string `json:"id"`
	PID      string `json:"pid"`
	Category string `json:"category"`
}

type UserLookupResultStruct struct {
	PID         uuid.UUID                  `json:"pid"`
	MessageType string                     `json:"type"`
	OrgType     string                     `json:"orgType"`
	Result      UserLookupResultDataStruct `json:"result"`
}

type UserLookupResultDataStruct struct {
	ID                      string   `json:"_id"`
	UID                     string   `json:"id"`
	PilotNames              []string `json:"pilotNames"`
	LoginTimes              []int64  `json:"loginTimes"`
	LogoutTimes             []int64  `json:"logoutTimes"`
	Kills                   int
	Deaths                  int
	Spawns                  SpawnTypeStruct          `json:"spawns"`
	ELO                     float64                  `json:"elo"`
	ELOHistory              []ELOHistoryStruct       `json:"eloHistory"`
	Rank                    int                      `json:"rank"`
	History                 string                   `json:"history"`
	DiscordID               string                   `json:"discordId"`
	IsBanned                bool                     `json:"isBanned"`
	TeamKills               int                      `json:"teamKills"`
	IgnoreKillsAgainstUsers []string                 `json:"ignoreKillsAgainstUsers"`
	EndOfSeasonStats        []EndOfSeasonStatsStruct `json:"endOfSeasonStats"`
	ELOFreeze               bool                     `json:"eloFreeze"`
}

type EndOfSeasonStatsStruct struct {
	Season    int     `json:"season"`
	Rank      int     `json:"rank"`
	ELO       float64 `json:"elo"`
	TeamKills int     `json:"teamKills"`
	History   string  `json:"history"`
}

type ELOHistoryStruct struct {
	Time int64   `json:"time"`
	ELO  float64 `json:"elo"`
}

type SpawnTypeStruct struct {
	EF24G   int `json:"0"`
	T55     int `json:"1"`
	Invalid int `json:"2"`
	AH94    int `json:"3"`
	F45A    int `json:"4"`
	FA26B   int `json:"5"`
	AV42C   int `json:"6"`
}

var Success = make(chan bool)

var Reconnecting = false

var SteamID64 string

var UserOnline bool

var FrontendWindow fyne.Window

var OnlineUsers OnlineStruct

var StopRichPresence = make(chan bool)

type UserStatsStruct struct {
	Kills              int
	Deaths             int
	Ratio              float64
	ELO                int
	LastSpawnTimestamp time.Time
	CurrentVehicle     string
	SpawnedIn          bool
}

var LatestUserStats UserStatsStruct

var WsStreamClosed = make(chan bool)

var Socket *gws.Conn
