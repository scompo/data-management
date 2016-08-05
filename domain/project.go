package domain

import "time"

type Project struct {
	Name         string
	Description  string
	CreationDate time.Time
	LastEdit     time.Time
}
