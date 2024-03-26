package api

import (
	"context"
	"encoding/base64"
	"net/http"
	"os"

	"github.com/henomis/restclientgo"
)

const (
	langfuseDefaultEndpoint = "https://cloud.langfuse.com"
)

type Client struct {
	restClient *restclientgo.RestClient
}

func New() *Client {
	langfuseHost := os.Getenv("LANGFUSE_HOST")
	if langfuseHost == "" {
		langfuseHost = langfuseDefaultEndpoint
	}

	publicKey := os.Getenv("LANGFUSE_PUBLIC_KEY")
	secretKey := os.Getenv("LANGFUSE_SECRET_KEY")

	restClient := restclientgo.New(langfuseHost)
	restClient.SetRequestModifier(func(req *http.Request) *http.Request {
		req.Header.Set("Authorization", basicAuth(publicKey, secretKey))
		return req
	})

	return &Client{
		restClient: restClient,
	}
}

func (c *Client) Ingestion(ctx context.Context, req *Ingestion, res *IngestionResponse) error {
	return c.restClient.Post(ctx, req, res)
}

func basicAuth(publicKey, secretKey string) string {
	auth := publicKey + ":" + secretKey
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
