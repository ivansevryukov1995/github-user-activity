package model

import (
	"strings"
)

func ConvertPayload(event string, e Event) Payload {
	switch event {
	case "Create":
		return Payload{RefType: e.Payload.RefType, Action: event + "d"}
	case "Delete":
		return Payload{RefType: e.Payload.RefType, Action: event + "d"}
	case "Gollum":
		return Payload{RefType: "page", Action: event + "ed"}
	case "Push":
		return Payload{RefType: "commit", Action: event + "ed"}
	case "Watch":
		return Payload{RefType: "repository", Action: event + "ed"}
	case "Issues":
		refType := "issue"
		return Payload{RefType: refType, Action: strings.Title(e.Payload.Action)}
	default:
		return Payload{RefType: strings.ToLower(event), Action: strings.Title(e.Payload.Action)}
	}
}
