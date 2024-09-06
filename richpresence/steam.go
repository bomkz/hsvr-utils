package richpresence

import (
	"log"

	"github.com/mrwaggel/gosteamconv"
	"golang.org/x/sys/windows/registry"
)

func findCurrentUID() (ID int32) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam\ActiveProcess`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer func(k registry.Key) {
		err := k.Close()
		if err != nil {

		}
	}(k)

	ID32, _, err := k.GetIntegerValue("ActiveUser")

	if err != nil {
		log.Fatal(err)
		return
	}

	return int32(ID32)
}

func convertID3ToID64(ID3 int32) int64 {
	IDString, err := gosteamconv.SteamInt32ToString(ID3)
	if err != nil {
		log.Println(err)
	}
	ID64, err := gosteamconv.SteamStringToInt64(IDString)
	if err != nil {
		log.Println(err)
	}
	return ID64
}
