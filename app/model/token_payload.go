package model

type TokenPayload struct {
	MessageType string `json:"messageType"`
	Token       string `json:"token"`
}
