package speechtotext

import (
	"context"
	"fmt"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

func GetAnswer(apiKey, filename string) (string, error) {
	if apiKey == "" || filename == "" {
		panic("API or filename value is empty")
	}

	c := openai.NewClient(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
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
