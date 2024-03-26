package api

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/henomis/langfuse-go/model"
)

const (
	ContentTypeJSON = "application/json"
)

type Request struct{}

type Ingestion struct {
	Batch []model.IngestionEvent `json:"batch"`
}

func (t *Ingestion) Path() (string, error) {
	return "/api/public/ingestion", nil
}

func (t *Ingestion) Encode() (io.Reader, error) {
	jsonBytes, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonBytes), nil
}

func (t *Ingestion) ContentType() string {
	return ContentTypeJSON
}
