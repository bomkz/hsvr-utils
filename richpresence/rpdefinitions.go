package richpresence

import (
	"time"

	"github.com/google/uuid"
	"github.com/lxzan/gws"
)

type subscribeStruct struct {
	MessageType string   `json:"type"`
	Data        []string `json:"data"`
}

type onlineStruct struct {
	PID         uuid.UUID          `json:"pid"`
	MessageType string             `json:"type"`
	Data        []onlineDataStruct `json:"data"`
}

type onlineDataStruct struct {
	Name string `json:"name"`
	UID  string `json:"id"`
	Team string `json:"team"`
}

type logStruct struct {
	PID         uuid.UUID     `json:"pid"`
	MessageType string        `json:"type"`
	Data        logDataStruct `json:"data"`
}

type logDataStruct struct {
	UserID    string `json:"userId"`
	PilotName string `json:"pilotName,omitempty"`
}

type spawnStruct struct {
	PID         uuid.UUID       `json:"pid"`
	MessageType string          `json:"type"`
	Data        spawnDataStruct `json:"data"`
}

type spawnDataStruct struct {
	User       userStruct       `json:"user"`
	ServerInfo serverInfoStruct `json:"serverInfo"`
}

type serverInfoStruct struct {
	OnlineUsers []string          `json:"onlineUsers"`
	MissionID   string            `json:"missionId"`
	Environment environmentStruct `json:"environment"`
}

type environmentStruct struct {
	TimeOfDay float64    `json:"tod"`
	Weather   int        `json:"weather"`
	Wind      windStruct `json:"wind"`
}

type windStruct struct {
	Heading   float64 `json:"heading"`
	Magnitude float64 `json:"mag"`
	Variance  float64 `json:"var"`
	Gusts     float64 `json:"gust"`
}

type userStruct struct {
	OwnerID   string   `json:"ownerId"`
	Occupants []string `json:"occupants"`
	Position  xyz      `json:"position"`
	Velocity  xyz      `json:"velocity"`
	Team      string   `json:"team"`
	Type      string   `json:"type"`
}

type xyz struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type killStruct struct {
	PID         uuid.UUID      `json:"pid"`
	MessageType string         `json:"type"`
	Data        killDataStruct `json:"data"`
}
type killDataStruct struct {
	Victim                  userStruct       `json:"victim"`
	Killer                  userStruct       `json:"killer"`
	ServerInfo              serverInfoStruct `json:"serverInfo"`
	Weapon                  string           `json:"weapon"`
	WeaponUUID              string           `json:"WeaponUuid"`
	PreviousDamagedByUserId string           `json:"previousDamagedByUserId"`
	PreviousDamagedByWeapon string           `json:"previousDamagedByWeapon"`
}

type deathStruct struct {
	PID         string          `json:"pid"`
	MessageType string          `json:"type"`
	Data        deathDataStruct `json:"data"`
}

type deathDataStruct struct {
	Victim     userStruct       `json:"victim"`
	ServerInfo serverInfoStruct `json:"serverInfo"`
}

type userStatsStruct struct {
	Kills              int
	Deaths             int
	Ratio              float64
	ELO                int
	LastSpawnTimestamp time.Time
	CurrentVehicle     string
	CurrentRank        int
	SpawnedIn          bool
}

var reconnectingRP bool

var attemptedAt time.Time

var success = make(chan bool)

var rpSuccess = make(chan bool)

var reconnecting = false

var steamID64 string

var userOnline bool

var latestUserStats userStatsStruct

var wsStreamClosed = make(chan bool)

var localSocket *gws.Conn

var onlineUsers onlineStruct
