package lubelogger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Construct a full request URL by combining the supplied endpoint with the
// api_url configuration value.
func endpointURL(endpoint string) string {
	return fmt.Sprintf("%s/%s", apiURI, endpoint)
}

// Value for the authorization header expected by th LubeLogger API.
func authorizationHeader() string {
	return authorization
}

// LubeLogger endpoint drop-in replacement for http.Get().
func GetEndpointWithContext(ctx context.Context, endpoint string) (*http.Response, error) {
	requestURL := endpointURL(endpoint)

	logger.DebugContext(ctx, "GetEndpoint called",
		"endpoint", endpoint,
		"url", requestURL,
	)

	apiRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("GetEndpointWithContext NewRequest: %w", err)
	}

	apiRequest.Header.Add("Authorization", authorizationHeader())

	apiResponse, err := http.DefaultClient.Do(apiRequest)

	logger.DebugContext(ctx, "GetEndpointWithContext apiResponse",
		"endpoint", endpoint,
		"statusCode", apiResponse.StatusCode,
		"proto", apiResponse.Proto,
		"contentLength", apiResponse.ContentLength,
	)

	return apiResponse, err
}

// Wrapped Get for standardized call of API GET endpoints.
func APIGet(endpoint string) ([]byte, error) {
	var ctx context.Context

	apiResponse, err := GetEndpointWithContext(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("APIGet GetEndpointWithContext: %w", err)
	}
	defer apiResponse.Body.Close()

	responseBody, err := io.ReadAll(apiResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("reading responseBody: %w", err)
	}

	logger.DebugContext(ctx, "LubeLogger APIGet",
		"responseBytes", len(responseBody),
	)

	return responseBody, nil
}

// LubeLogger endpoint drop-in replacement for http.PostForm().
func PostFormEndpointWithContext(ctx context.Context, endpoint string, data url.Values) (*http.Response, error) {
	requestURL := endpointURL(endpoint)

	logger.DebugContext(ctx, "PostFormEndpoint called",
		"endpoint", endpoint,
		"url", requestURL,
		"data", data,
	)

	requestBody := strings.NewReader(data.Encode())

	apiRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, requestBody)
	if err != nil {
		return nil, fmt.Errorf("PostFormEndpointWithContext NewRequest: %w", err)
	}

	apiRequest.Header.Add("Authorization", authorizationHeader())
	apiRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	apiResponse, err := http.DefaultClient.Do(apiRequest)

	logger.DebugContext(ctx, "PostFormEndpoint apiResponse",
		"endpoint", endpoint,
		"statusCode", apiResponse.StatusCode,
		"proto", apiResponse.Proto,
		"contentLength", apiResponse.ContentLength,
	)

	return apiResponse, err
}

// Wrapped PostForm for standardized call of API POST endpoints.
func APIPostForm(endpoint string, data url.Values) (PostResponse, error) {
	var (
		ctx      context.Context
		response PostResponse
	)

	apiResponse, err := PostFormEndpointWithContext(ctx, endpoint, data)
	if err != nil {
		return PostResponse{}, fmt.Errorf("APIPostForm PostFormEndpoint: %w", err)
	}
	defer apiResponse.Body.Close()

	responseBody, err := io.ReadAll(apiResponse.Body)
	if err != nil {
		return PostResponse{}, fmt.Errorf("APIPostForm reading responseBody: %w", err)
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return PostResponse{}, fmt.Errorf("unmarshalling json: %w", err)
	}

	logger.Debug("LubeLogger APIPostForm",
		"success", response.Success,
		"message", response.Message,
		"status", apiResponse.StatusCode,
		"formBytes", len(data.Encode()),
		"formFields", len(data),
	)

	if !response.Success {
		return response, fmt.Errorf("post: %s", response.Message)
	}

	return response, nil
}
