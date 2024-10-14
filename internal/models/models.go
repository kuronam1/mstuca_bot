package models

type User struct {
	ID        int64
	Name      string
	GroupName string
	GroupID   int
	Subgroup  int
	State     int
	Title     string
}

type Group struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type Filter struct {
}

type Schedule struct {
}
