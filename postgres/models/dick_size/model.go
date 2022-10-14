package models

import "time"

type DickSize struct {
	Id           int       `json:"id"`
	UsedId       int       `json:"userId"`
	Fname        string    `json:"fname"`
	Lname        string    `json:"lname"`
	Username     string    `json:"username"`
	Dick_size    int8      `json:"dick_size"`
	Measure_date time.Time `json:"measure_date"`
	Chat_id      int       `json:"chat_id"`
	Is_group     bool      `json:"is_group"`
}

type UserCredentials struct {
	UserId   int64  `json:"user_id"`
	Fname    string `json:"fname"`
	Username string `json:"username"`
	Lname    string `json:"lname"`
}

type Duel struct {
	CallerUserId int64 `json:"caller_user_id"`
	CalledUserId int64 `json:"called_user_id"`
	ChatID       int64 `json:"chat_id"`
	Bet          int   `json:"bet"`
	Winner       int64 `json:"winner"`
}
