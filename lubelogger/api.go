package lubelogger

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Get(endpoint string) ([]byte, error) {
	uri := fmt.Sprintf("%s/%s", api_uri, endpoint)
	log.WithFields(log.Fields{
		"uri": uri,
	}).Trace("API URI")

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	req.Header.Add("Authorization", authorization)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}

	return body, nil
}

func Post(endpoint string, data url.Values) (response PostResponse, err error) {
	uri := fmt.Sprintf("%s/%s", api_uri, endpoint)

	log.WithFields(log.Fields{
		"uri": uri,
	}).Trace("API URI")

	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(data.Encode()))
	if err != nil {
		return PostResponse{}, fmt.Errorf("building request: %w", err)
	}

	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return PostResponse{}, fmt.Errorf("sending request: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PostResponse{}, fmt.Errorf("reading body: %w", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return PostResponse{}, fmt.Errorf("unmarshalling json: %w", err)
	}

	log.WithFields(log.Fields{
		"uri":        uri,
		"success":    response.Success,
		"message":    response.Message,
		"status":     resp.StatusCode,
		"formBytes":  len(data.Encode()),
		"formFields": len(data),
	}).Debug("LubeLogger API Post")

	if !response.Success {
		return response, fmt.Errorf("post: %s", response.Message)
	}

	return response, nil
}
