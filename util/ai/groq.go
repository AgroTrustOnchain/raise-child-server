package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"raise-child/constants/env"
	"strings"
)

type GroqClient struct{}

const defaultGroqAPIURL string = "https://api.groq.com/openai/v1/chat/completions"

type AiResponse struct {
	Result string `json:"result"`
	Reason string `json:"reason"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

// TextPart creates a text ContentPart.
func TextPart(text string) ContentPart {
	return ContentPart{Type: "text", Text: text}
}

// ImagePart creates an image_url ContentPart.
func ImagePart(url string) ContentPart {
	return ContentPart{Type: "image_url", ImageURL: &ImageURL{URL: url}}
}

type Message struct {
	Role    string        `json:"role"`
	Content []ContentPart `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChoiceMessage holds the content returned by the model for a single choice.
type ChoiceMessage struct {
	Content string `json:"content"`
}

// Choice represents one completion choice from the Groq API.
type Choice struct {
	Message ChoiceMessage `json:"message"`
}

// ApiError is the error object returned by the Groq API on failure.
type ApiError struct {
	Message string `json:"message"`
}

// ChatResponse is the top-level response envelope from the Groq API.
// Choices is a slice because the API always returns an array.
type ChatResponse struct {
	Choices []Choice  `json:"choices"`
	Error   *ApiError `json:"error,omitempty"`
}

// LabeledImage represents an image URL with a description label.
type LabeledImage struct {
	URL   string // Full image URL (http/https)
	Label string // Description of what the image is (e.g., "Child's birth certificate")
}

// AskWithImages sends a prompt together with multiple image URLs to the Groq
// vision model and returns the raw text content of the first choice.
func (g *GroqClient) AskWithImages(ctx context.Context, apiKey, prompt string, imageURLs []string) (string, error) {
	groqAPIURL := os.Getenv(env.GROQ_API_URL)
	if groqAPIURL == "" {
		groqAPIURL = defaultGroqAPIURL
	}

	model := os.Getenv(env.GROQ_MODEL)
	if model == "" {
		return "", fmt.Errorf("missing environment variable %s", env.GROQ_MODEL)
	}

	// Validate image URLs have proper schemes
	for i, url := range imageURLs {
		if url == "" {
			return "", fmt.Errorf("image URL at index %d is empty", i)
		}
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			return "", fmt.Errorf("image URL at index %d does not have http or https scheme: %s", i, url)
		}
	}

	parts := []ContentPart{TextPart(prompt)}
	for _, url := range imageURLs {
		parts = append(parts, ImagePart(url))
	}

	reqBody := ChatRequest{
		Model: model,
		Messages: []Message{
			{Role: "user", Content: parts},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", groqAPIURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call API: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBytes, &chatResp); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("groq API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 || chatResp.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("no content returned")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// AskWithDescribedImages sends a prompt together with labeled images to the Groq vision model.
// Each image is preceded by a text description so the model knows what each image represents.
func (g *GroqClient) AskWithDescribedImages(ctx context.Context, apiKey, prompt string, labeledImages []LabeledImage) (string, error) {
	groqAPIURL := os.Getenv(env.GROQ_API_URL)
	if groqAPIURL == "" {
		groqAPIURL = defaultGroqAPIURL
	}

	model := os.Getenv(env.GROQ_MODEL)
	if model == "" {
		return "", fmt.Errorf("missing environment variable %s", env.GROQ_MODEL)
	}

	// Validate images
	for i, img := range labeledImages {
		if img.URL == "" {
			return "", fmt.Errorf("image URL at index %d is empty", i)
		}
		if !strings.HasPrefix(img.URL, "http://") && !strings.HasPrefix(img.URL, "https://") {
			return "", fmt.Errorf("image URL at index %d does not have http or https scheme: %s", i, img.URL)
		}
	}

	// Build content parts with descriptions before each image
	parts := []ContentPart{TextPart(prompt)}
	for _, img := range labeledImages {
		if img.Label != "" {
			parts = append(parts, TextPart(img.Label))
		}
		parts = append(parts, ImagePart(img.URL))
	}

	reqBody := ChatRequest{
		Model: model,
		Messages: []Message{
			{Role: "user", Content: parts},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", groqAPIURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call API: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBytes, &chatResp); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("groq API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 || chatResp.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("no content returned")
	}

	return chatResp.Choices[0].Message.Content, nil
}
