package model

import "time"

type Url struct {
	Hash      string    `json:"hash"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}
