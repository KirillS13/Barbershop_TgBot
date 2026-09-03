package models

import "time"

type Master struct {
	Id        int
	Name      string
	Phone     string
	IsActive  bool
	CreatedAt time.Time
}
