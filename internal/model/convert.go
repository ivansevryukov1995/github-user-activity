package model

import (
	"strings"
)

func ConvertPayloadToAction(event string, e Event) Action {
	switch event {
	case "Create":
		return Action{RefType: e.Payload.RefType, Action: event + "d"}
	case "Delete":
		return Action{RefType: e.Payload.RefType, Action: event + "d"}
	case "Gollum":
		return Action{RefType: "page", Action: event + "ed"}
	case "Push":
		return Action{RefType: "commit", Action: event + "ed"}
	case "Watch":
		return Action{RefType: "repository", Action: event + "ed"}
	case "Issues":
		refType := "issue"
		return Action{RefType: refType, Action: strings.Title(e.Payload.Action)}
	default:
		return Action{RefType: strings.ToLower(event), Action: strings.Title(e.Payload.Action)}
	}
}
