package domain

import (
	"encoding/json"
	"github.com/minio/minio-go"
	"io"
	"time"
)

var projects map[string]Project = make(map[string]Project)

type ProjectInfo struct {
	Description string
}

type Project minio.BucketInfo

var currentTime = time.Now

func (p *Project) Encode(w io.Writer) {
	p.CreationDate = currentTime()
	enc := json.NewEncoder(w)
	enc.Encode(p)
	return
}
