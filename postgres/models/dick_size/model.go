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
