package models

type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	ApiKey string `json:"api_key"`
}
