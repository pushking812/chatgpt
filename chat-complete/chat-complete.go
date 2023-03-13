package chatcomplete

import (
	"context"
	"fmt"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

func GetAnswer(apiKey, msg string) (string, error) {
	if apiKey == "" || msg == "" {
		panic("Missing API KEY or message empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}
