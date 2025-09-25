package structs

import "time"

type Marathon struct {
	Id                   string
	Name                 string
	Webhook              string
	StartDate            time.Time
	EndDate              time.Time
	SubmissionsStartDate time.Time
	SubmissionsEndDate   time.Time
}
