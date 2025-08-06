package pkg

import (
	"bytes"
	"net/http"
)

func Post(client *http.Client, url string, contentType string, data []byte) *http.Response {
	resp, err := client.Post(
		url,
		contentType,
		bytes.NewReader(data),
	)
	CheckError(err)

	if resp.StatusCode != http.StatusOK {
		// TODO: check status code, potentially wait and re-send
		panic("bootstrap failed: " + resp.Status)
	}

	return resp
}
