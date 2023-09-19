package artificialintelligence

import (
	"context"
	"errors"
	"github.com/sashabaranov/go-openai"
	"log/slog"
	"os"
)

var (
	errInvalidKey = errors.New("invalid openai API key")
)

type SuperAI interface {
	AddMessage(string)
	Ask(string) (string, error)
}

func NewClient() (SuperAI, error) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return nil, errInvalidKey
	}
	openaiCli := openai.NewClient(key)
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "you are a helpful chat-bot",
		},
	}

	return &superAI{openaiCli, messages}, nil
}

type openAI interface {
	CreateChatCompletion(context.Context, openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

type superAI struct {
	client   openAI
	messages []openai.ChatCompletionMessage
}

func (ai *superAI) AddMessage(msg string) {
	ai.messages = append(ai.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: msg,
	})
}

func (ai *superAI) Ask(message string) (string, error) {
	ai.AddMessage(message)

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: ai.messages,
	}

	resp, err := ai.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		slog.Error("chat completion error", "err", err)
		return "", err
	}

	answer := resp.Choices[0].Message
	ai.messages = append(ai.messages, answer)

	return answer.Content, nil
}
