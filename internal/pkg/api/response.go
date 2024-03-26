package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/henomis/restclientgo"
)

type Response struct {
	Code      int       `json:"-"`
	RawBody   *string   `json:"-"`
	Successes []Success `json:"successes"`
	Errors    []Error   `json:"errors"`
}

type Success struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
}

type Error struct {
	ID      string `json:"id"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func (r *Response) IsSuccess() bool {
	return r.Code < http.StatusBadRequest
}

func (r *Response) SetStatusCode(code int) error {
	r.Code = code
	return nil
}

func (r *Response) SetBody(body io.Reader) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	s := string(b)
	r.RawBody = &s

	return nil
}

func (r *Response) AcceptContentType() string {
	return ContentTypeJSON
}

func (r *Response) Decode(body io.Reader) error {
	return json.NewDecoder(body).Decode(r)
}

func (r *Response) SetHeaders(_ restclientgo.Headers) error {
	return nil
}

type IngestionResponse struct {
	Response
}
