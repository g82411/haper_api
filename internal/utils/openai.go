package utils

import (
	"bytes"
	"encoding/json"
	"hyper_api/internal/config"
	"net/http"
)

const (
	IMAGE_GENERATE_URL = "https://api.openai.com/v1/images/generations"
)

func sendPostRequest(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	setting := config.GetConfig()
	key := setting.OpenAIKey
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GeneratePhotoUsingDallE3(prompt string, n int) ([]string, error) {
	type ImageGenerateRequest struct {
		Prompt string `json:"prompt"`
		N      int    `json:"n"`
		Size   string `json:"size"`
		Model  string `json:"model"`
	}
	type Image struct {
		RevisedPrompt string `json:"revised_prompt"`
		URL           string `json:"url"`
	}
	type ApiResponse struct {
		Created int64   `json:"created"`
		Data    []Image `json:"data"`
	}
	type ApiError struct {
		Code    *string `json:"code"`
		Message string  `json:"message"`
		Param   *string `json:"param"`
		Type    string  `json:"type"`
	}
	type ErrorResponse struct {
		Error ApiError `json:"error"`
	}

	ans := make([]string, n)
	body := ImageGenerateRequest{
		Prompt: prompt,
		N:      n,
		Size:   "1024x1024",
		Model:  "dall-e-3",
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return ans, err
	}
	resp, err := sendPostRequest(IMAGE_GENERATE_URL, jsonBody)
	if err != nil {
		return ans, err
	}
	defer resp.Body.Close()
	var apiResponse ApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return ans, err
	}
	if len(apiResponse.Data) == 0 {
		var errorResponse ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return ans, err
		}
	}
	for i, url := range apiResponse.Data {
		ans[i] = url.URL
	}
	return ans, nil
}
