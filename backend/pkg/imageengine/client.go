package imageengine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type Text2ImageRequest struct {
	Prompt         string  `json:"prompt"`
	ModelID        string  `json:"model_id,omitempty"`
	NegativePrompt string  `json:"negative_prompt,omitempty"`
	Width          int     `json:"width,omitempty"`
	Height         int     `json:"height,omitempty"`
	CFGScale       float64 `json:"cfg_scale,omitempty"`
	Steps          int     `json:"steps,omitempty"`
	Seed           int64   `json:"seed,omitempty"`
	Style          string  `json:"style,omitempty"`
	Quality        string  `json:"quality,omitempty"`
}

type ImageEditRequest struct {
	Image          string  `json:"image"`
	Prompt         string  `json:"prompt"`
	ModelID        string  `json:"model_id,omitempty"`
	NegativePrompt string  `json:"negative_prompt,omitempty"`
	Mask           string  `json:"mask,omitempty"`
	Strength       float64 `json:"strength,omitempty"`
	CFGScale       float64 `json:"cfg_scale,omitempty"`
	Steps          int     `json:"steps,omitempty"`
	Seed           int64   `json:"seed,omitempty"`
}

type GenerateResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Status  string   `json:"status"`
	TaskID  string   `json:"task_id"`
	Images  []string `json:"images"`
	Quota   int      `json:"quota"`
}

type AIModel struct {
	ID                   int                    `json:"id"`
	ModelID              string                 `json:"model_id"`
	ModelName            string                 `json:"model_name"`
	Provider             string                 `json:"provider"`
	Description          string                 `json:"description"`
	Capabilities         []string               `json:"capabilities"`
	SupportedParameters  []string               `json:"supported_parameters"`
	IsActive             bool                   `json:"is_active"`
	Priority             int                    `json:"priority"`
	Metadata             map[string]interface{} `json:"metadata"`
	ParameterConstraints map[string]interface{} `json:"parameter_constraints"`
	CreatedAt            string                 `json:"created_at"`
	UpdatedAt            string                 `json:"updated_at"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewClient() *Client {
	baseURL := viper.GetString("image_engine.base_url")
	apiKey := viper.GetString("image_engine.api_key")
	timeout := viper.GetInt("image_engine.timeout")
	if timeout == 0 {
		timeout = 120
	}

	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

func (c *Client) doRequest(method, path string, body interface{}) (*Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body failed: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.baseURL + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	var result Response
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return &result, nil
}

func (c *Client) Text2Image(req *Text2ImageRequest) (*GenerateResponse, error) {
	resp, err := c.doRequest("POST", "/api/v1/text2image", req)
	if err != nil {
		return nil, err
	}

	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data type")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var result GenerateResponse
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) ImageEdit(req *ImageEditRequest) (*GenerateResponse, error) {
	resp, err := c.doRequest("POST", "/api/v1/image/edit", req)
	if err != nil {
		return nil, err
	}

	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data type")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var result GenerateResponse
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetModels() ([]AIModel, error) {
	resp, err := c.doRequest("GET", "/api/v1/models", nil)
	if err != nil {
		return nil, err
	}

	data, ok := resp.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data type")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var models []AIModel
	if err := json.Unmarshal(jsonData, &models); err != nil {
		return nil, err
	}

	return models, nil
}

func (c *Client) GetModelsByCapability(capability string) ([]AIModel, error) {
	resp, err := c.doRequest("GET", "/api/v1/models/capability/"+capability, nil)
	if err != nil {
		return nil, err
	}

	data, ok := resp.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data type")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var models []AIModel
	if err := json.Unmarshal(jsonData, &models); err != nil {
		return nil, err
	}

	return models, nil
}

func (c *Client) HealthCheck() (bool, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/health", nil)
	if err != nil {
		return false, err
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
