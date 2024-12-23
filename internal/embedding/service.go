package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type EmbeddingService interface {
	GetEmbedding(text string) ([]float64, error)
}

// OllamaEmbedding implements embedding using local Ollama server
type OllamaEmbedding struct {
	baseURL string
	model   string
}

func NewOllamaEmbedding(baseURL, model string) *OllamaEmbedding {
	if baseURL == "" {
		baseURL = "http://localhost:11434" // Fallback default
	}
	if model == "" {
		model = "nomic-embed-text" // Fallback default
	}

	return &OllamaEmbedding{
		baseURL: baseURL,
		model:   model,
	}
}

type ollamaRequest struct {
	Model   string  `json:"model"`
	Prompt  string  `json:"prompt"`
	Options Options `json:"options"`
}

type Options struct {
	Temperature float64 `json:"temperature"`
}

type ollamaResponse struct {
	Embedding []float64 `json:"embedding"`
}

func (s *OllamaEmbedding) GetEmbedding(text string) ([]float64, error) {
	url := fmt.Sprintf("%s/api/embeddings", s.baseURL)

	requestBody := ollamaRequest{
		Model:  s.model,
		Prompt: text,
		Options: Options{
			Temperature: 0.0, // We want deterministic results for embeddings
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama API returned status code %d", resp.StatusCode)
	}

	var result ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding Ollama response: %w", err)
	}

	if len(result.Embedding) == 0 {
		return nil, fmt.Errorf("no embedding returned from Ollama")
	}

	return result.Embedding, nil
}
