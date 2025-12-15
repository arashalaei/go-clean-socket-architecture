package dto

import "encoding/json"

type Request struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
