package model

import "time"

type EventRequest struct {
	Domain 		string		`json:"domain"`
	Path 		string 		`json:"path"`
	Referrer	string		`json:"referrer"`
	ScreenSize 	string		`json:"screen_size"`
}

type PageView struct {
	ID			int64
	Domain		string
	Path    	string
	Referrer 	string
	CountryCode string
	ScreenSize string
	Browser string
	OS string
	VisitorHash string
	CreatedAt   time.Time
}