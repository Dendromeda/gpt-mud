package gpt

import (
	"bufio"
	"fmt"
	"os"
)

type Chat struct {
	model   string
	Entries []*chatEntry
	client  *Client
}

type chatEntry struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewChat(client *Client, model string) *Chat {
	return &Chat{
		model:  model,
		client: client,
	}
}

func (c *Chat) Chat(role, content string) (string, error) {
	c.Entries = append(c.Entries, &chatEntry{
		Role:    role,
		Content: content,
	})
	resp, err := c.client.Chat(c)
	if err != nil {
		return "", err
	}
	c.Entries = append(c.Entries, resp)
	return resp.Content, nil
}

func (c *Chat) Prompt() error {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return nil
	}
	output, err := c.Chat("user", input)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return nil
	}
	fmt.Println("GPT 4:", output)
	return nil
}
