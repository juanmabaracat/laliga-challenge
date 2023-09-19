package artificialintelligence

import (
	"context"
	"errors"
	"github.com/sashabaranov/go-openai"
	"os"
	"reflect"
	"testing"
)

type mockOpenAI struct {
	*openai.Client
	response openai.ChatCompletionResponse
	err      error
}

func (mock *mockOpenAI) CreateChatCompletion(context.Context,
	openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	return mock.response, mock.err
}

func Test_superAI_AddMessage(t *testing.T) {
	superAI := superAI{
		client:   openai.NewClient("test"),
		messages: nil,
	}

	if len(superAI.messages) != 0 {
		t.Errorf("messages should be empty: %v", superAI.messages)
	}

	firstQuestion := "this is a first question"
	superAI.AddMessage(firstQuestion)

	if len(superAI.messages) != 1 {
		t.Errorf("it should have 1 message: %v", superAI.messages)
	}

	if !reflect.DeepEqual(superAI.messages[0].Content, firstQuestion) {
		t.Errorf("GOT=%s, WANT=%s", superAI.messages[0].Content, firstQuestion)
	}
}

func Test_superAI_Ask(t *testing.T) {
	tests := []struct {
		name         string
		question     string
		expectedResp string
		expectedErr  error
	}{
		{
			"should respond correctly",
			"Where was Lionel Messi born?",
			"Rosario, Santa Fe, Argentina",
			nil,
		},
		{
			"should return error when AI doesn't work",
			"Can you help me",
			"",
			errors.New("connection error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ai := &superAI{
				client:   createMock(tt.expectedResp, tt.expectedErr),
				messages: nil,
			}

			got, err := ai.Ask(tt.question)
			if (tt.expectedErr != nil) && err != tt.expectedErr {
				t.Errorf("Ask() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if got != tt.expectedResp {
				t.Errorf("Ask() got = %v, want %v", got, tt.expectedResp)
			}
		})
	}
}

func createMock(msg string, err error) *mockOpenAI {
	resp := openai.ChatCompletionResponse{}
	resp.Choices = []openai.ChatCompletionChoice{
		{
			Index:   0,
			Message: openai.ChatCompletionMessage{Content: msg},
		},
	}
	return &mockOpenAI{
		response: resp,
		err:      err,
	}
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		wantErr error
	}{
		{
			name:    "should create a client without error",
			apiKey:  "test-key",
			wantErr: nil,
		},
		{
			name:    "should return an error when the API key is not present",
			apiKey:  "",
			wantErr: errInvalidKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("OPENAI_API_KEY", tt.apiKey)
			_, err := NewClient()
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("GOT=%v, EXPECTED=%v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr == nil {
				t.Errorf("GOT=%v, EXPECTED not error", err)
			}

		})
	}
}
