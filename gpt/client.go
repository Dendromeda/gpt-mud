package gpt

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	APIKey   string
	Endpoint string
	client   *resty.Client
}

type chatCompletion struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

func NewClient(apiKey, endpoint string) *Client {
	return &Client{
		APIKey:   apiKey,
		Endpoint: endpoint,
		client:   resty.New(),
	}
}

func (c *Client) Chat(chat *Chat) (*chatEntry, error) {

	response, err := c.client.R().
		SetAuthToken(c.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":      chat.model,
			"messages":   chat.Entries,
			"max_tokens": 50,
		}).
		Post(c.Endpoint)
	if err != nil {
		return nil, err
	}

	body := response.Body()

	var data *chatCompletion
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error while decoding JSON response:", err)
		return nil, err
	}
	if len(data.Choices) == 0 {
		fmt.Println(string(body))
		return nil, fmt.Errorf("no chat choices returned")

	}

	chatEntry := &chatEntry{
		Role:    data.Choices[0].Message.Role,
		Content: data.Choices[0].Message.Content,
	}

	return chatEntry, nil
}
