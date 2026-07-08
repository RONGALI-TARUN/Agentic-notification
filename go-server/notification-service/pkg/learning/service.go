package learning

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Service struct {
	baseURL string
}

func NewService(baseURL string) *Service {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &Service{baseURL: baseURL}
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func (s *Service) Summarize(text string) (string, error) {
	if text == "" {
		return "", fmt.Errorf("text required")
	}

	payload := OllamaRequest{
		Model:  "llama3.2",
		Prompt: fmt.Sprintf("Summarize this feedback into one short recommendation:\n\n%s", text),
		Stream: false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(s.baseURL+"/api/generate", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama request failed with status %d", resp.StatusCode)
	}

	var out OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.Response, nil
}
