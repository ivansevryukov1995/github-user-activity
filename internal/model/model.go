package model

import "time"

type Actor struct {
	ID           int    `json:"id"`
	Login        string `json:"login"`
	DisplayLogin string `json:"display_login"`
	GravatarID   string `json:"gravatar_id"`
	URL          string `json:"url"`
	AvatarURL    string `json:"avatar_url"`
}

type Repo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Payload struct {
	Ref     string `json:"ref,omitempty"`
	RefType string `json:"ref_type,omitempty"`
	Action  string `json:"action,omitempty"`
}

type Event struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Actor     Actor     `json:"actor"`
	Repo      Repo      `json:"repo"`
	Payload   Payload   `json:"payload,omitempty"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	NameUser       string    `json:"name_user"`
	CreatedAtCache time.Time `json:"created_at_cache"`
	Events         []Event
}
