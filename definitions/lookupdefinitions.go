package definitions

import "github.com/google/uuid"

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

type EndOfSeasonStatsStruct struct {
	Season    int     `json:"season"`
	Rank      int     `json:"rank"`
	ELO       float64 `json:"elo"`
	TeamKills int     `json:"teamKills"`
	History   string  `json:"history"`
}
