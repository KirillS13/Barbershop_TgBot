package models

import "time"

type Appointment struct {
	Id        int
	ClientId  int64
	ServerId  int
	MasterId  int
	Time      time.Time
	Status    string
	CreatedAt time.Time
}
