package model

type TokenPayload struct {
	MessageType string `json:"mt"`
	Token       string `json:"t"`
}
