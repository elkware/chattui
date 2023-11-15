package internal

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"chattui/internal/config"
)

type ResponseMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type ResponseChoice struct {
	FinishReason string          `json:"finish_reason"`
	Index        int             `json:"index"`
	Message      ResponseMessage `json:"message"`
}

type ChatGPTResponse struct {
	Choices []ResponseChoice `json:"choices"`
	Created int              `json:"created"`
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Model   string           `json:"model"`
	Usage   struct {
		CompletionTokens int `json:"completion_tokens"`
		PromptTokens     int `json:"prompt_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type RequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPTRequest struct {
	Model    string           `json:"model"`
	Messages []RequestMessage `json:"messages"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Chat struct {
	Name     string        `json:"name"`
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

var Chats = make([]*Chat, 0)

func CallChatGPT(chat *Chat) (string, error) {
	chatGptMessages := make([]RequestMessage, 0)
	if config.AppConfig.CustomInstructions != "" {
		chatGptMessages = append(chatGptMessages, RequestMessage{
			Role:    "system",
			Content: config.AppConfig.CustomInstructions,
		})
	}
	for _, message := range chat.Messages {
		chatGptMessages = append(chatGptMessages, RequestMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	r := ChatGPTRequest{
		Model:    chat.Model,
		Messages: chatGptMessages,
	}
	jso, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", strings.NewReader(string(jso)))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.AppConfig.ApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	var result ChatGPTResponse
	responseBody, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		panic(err)
	}
	if len(result.Choices) == 0 {
		return "", errors.New("no response from GPT")
	}
	return result.Choices[0].Message.Content, nil
}
