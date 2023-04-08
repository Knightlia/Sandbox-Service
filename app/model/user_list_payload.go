package model

type UserListPayload struct {
	MessageType string   `json:"mt"`
	UserList    []string `json:"ul"`
}
