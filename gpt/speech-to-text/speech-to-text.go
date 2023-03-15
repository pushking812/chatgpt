package speechtotext

import (
	"context"
	"fmt"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

func GetAnswer(apiKey, filename, t string) (string, error) {
	if apiKey == "" || filename == "" {
		panic("API or filename value is empty")
	}

	c := openai.NewClient(apiKey)

	timeout, err := time.ParseDuration(t)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: filename,
	}
	resp, err := c.CreateTranscription(ctx, req)
	if err != nil {
		return "", fmt.Errorf("transcription error: %v", err)
	}

	return resp.Text, nil
}
