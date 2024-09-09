package apifrontend

import (
	"bytes"
	"encoding/json"

	"github.com/bomkz/vtolvr-utils/definitions"
)

func handleWs(message bytes.Buffer) {

	var newLookupType definitions.LookupResultsTypeStruct

	json.Unmarshal(message.Bytes(), &newLookupType)

	var newMessage messageStruct
	newMessage.message = message
	newMessage.PID = newLookupType.PID

	messages = append(messages, newMessage)

}
