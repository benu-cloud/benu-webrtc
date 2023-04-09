package message

import (
	"encoding/json"
)

type MessageType string

const (
	SessionDescriptionMessage MessageType = "sdp"
	IceCandidateMessage       MessageType = "ice"
)

type GenericMessage struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type PayloadTarget string

const (
	Video PayloadTarget = "video"
	Audio PayloadTarget = "audio"
)

type GenericPayload interface{}

type SessionDescriptionPayload struct {
	From               string          `json:"from"`
	Target             PayloadTarget   `json:"target"`
	SessionDescription json.RawMessage `json:"content"`
}

type IceCandidatePayload struct {
	From         string          `json:"from"`
	Target       PayloadTarget   `json:"target"`
	IceCandidate json.RawMessage `json:"content"`
}
