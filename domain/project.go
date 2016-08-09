package domain

import (
	"encoding/json"
	"io"
	"time"
)

var projects map[string]Project = make(map[string]Project)

type Project struct {
	Name         string
	Description  string
	CreationDate time.Time
	LastEdit     time.Time
}

var currentTime = time.Now

func (p *Project) Encode(w io.Writer) {
	p.CreationDate = currentTime()
	p.LastEdit = currentTime()
	enc := json.NewEncoder(w)
	enc.Encode(p)
	return
}
