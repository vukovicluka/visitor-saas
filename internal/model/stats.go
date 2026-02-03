package model

type SummaryStats struct {
	TotalViews		int		  		`json:"total_views"`
	UniqueVisitors	int				`json:"unique_visitors"`
	ViewsPerDay		[]DailyStat		`json:"views_per_day"`
}

type DailyStat struct {
	Date		string		`json:"date"`
	Views		int			`json:"views"`
	Visitors	int			`json:"visitors"`
}

type PageStats struct {
	Path		string		`json:"path"`
	Views		int			`json:"views"`
	Visitors	int			`json:"visitors"`
}

type ReferrerStats struct {
	Referrer	string		`json:"referrer"`
	Views		int			`json:"views"`
	Visitors	int			`json:"visitors"`
}