package model

type MessagePayload struct {
	MessageType string `json:"mt"`
	Sender      string `json:"s"`
	Message     string `json:"m"`
	Timestamp   int64  `json:"t"`
}
