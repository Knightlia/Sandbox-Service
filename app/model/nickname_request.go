package model

import (
	"errors"
	"net/http"
	"strings"
)

type NicknameRequest struct {
	Nickname string
}

func (r NicknameRequest) Bind(_ *http.Request) error {
	if len(r.Nickname) == 0 {
		return errors.New("error.nickname.required")
	}
	if strings.Contains(r.Nickname, " ") {
		return errors.New("error.nickname.format")
	}
	return nil
}
