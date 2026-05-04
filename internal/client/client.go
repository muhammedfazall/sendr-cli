package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	token      string
	httpClient *http.Client
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}

func New(baseURL, token, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) do(method, path string, body any, useAPIKey bool) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("encode request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if useAPIKey {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	} else {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(data, &apiErr); err != nil {
			return nil, resp.StatusCode, fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
		return nil, resp.StatusCode, &apiErr
	}

	return data, resp.StatusCode, nil
}

// Auth

func (c *Client) GetToken(path string) (string, error) {
	data, _, err := c.do("GET", path, nil, false)
	if err != nil {
		return "", err
	}
	var result struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("parse token response: %w", err)
	}
	return result.Token, nil
}

func (c *Client) Logout() error {
	_, _, err := c.do("POST", "/auth/logout", nil, false)
	return err
}

// Keys

type APIKey struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Prefix    string `json:"prefix"`
	CreatedAt string `json:"created_at"`
}

type CreateKeyResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
	APIKey string `json:"api_key"`
}

func (c *Client) CreateKey(name string) (*CreateKeyResponse, error) {
	data, _, err := c.do("POST", "/apikeys", map[string]string{"name": name}, false)
	if err != nil {
		return nil, err
	}
	var result CreateKeyResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return &result, nil
}

func (c *Client) ListKeys() ([]APIKey, error) {
	data, _, err := c.do("GET", "/apikeys", nil, false)
	if err != nil {
		return nil, err
	}
	var keys []APIKey
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return keys, nil
}

func (c *Client) RevokeKey(id string) error {
	_, _, err := c.do("DELETE", "/apikeys/"+id, nil, false)
	return err
}

// Email

type SendEmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type SendEmailResponse struct {
	JobID   string `json:"job_id"`
	Message string `json:"message"`
}

func (c *Client) SendEmail(to, subject, body string) (*SendEmailResponse, error) {
	data, _, err := c.do("POST", "/emails/send", SendEmailRequest{
		To:      to,
		Subject: subject,
		Body:    body,
	}, true)
	if err != nil {
		return nil, err
	}
	var result SendEmailResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return &result, nil
}

func (c *Client) GetJob(id string) (map[string]any, error) {
	data, _, err := c.do("GET", "/emails/"+id, nil, true)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return result, nil
}