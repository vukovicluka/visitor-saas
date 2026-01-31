package model

import "time"

type EventRequest struct {
	Domain 		string		`json:"domain"`
	Path 		string 		`json:"path"`
	Referrer	string		`json:"referrer"`
}

type PageView struct {
	ID			int64
	Domain		string
	Path    	string
	Referrer 	string
	CountryCode string
	VisitorHash string
	CreatedAt   time.Time
}