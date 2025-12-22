package chat

import (
	"context"
	"errors"
	"fmt"
	"time"

	pp "github.com/Mushroom-MSL1L/UltimateDesktopPet/pkg/print"
	"google.golang.org/genai"
)

func newGeminiClient(apiKey string) (*genai.Client, error) {
	ctx := context.Background()
	pp.Assert(pp.Chat, "newGeminiClient: apiKey length %d", len(apiKey))
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		err = fmt.Errorf(
			"newGeminiClient: \n\tfailed to load .env file: %w \n\n\t"+
				"We need user to provide a valid genimi api at %s, content like \"GEMINI_API_KEY=<your_api_key>\"\n\t"+
				"You can generate a key by the link https://aistudio.google.com/app/api-keys",
			err, envFilePathForGemini)
		return nil, err
	}
	return client, nil
}

func (c *ChatMeta) chatWithModel(ctx context.Context, request string) (response string, err error) {
	if c == nil || c.Client == nil {
		return "", fmt.Errorf("genai client is not initialized")
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := c.Client.Models.GenerateContent(
		timeoutCtx,
		model,
		genai.Text(request),
		nil,
	)

	if err != nil {
		if errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) {
			return "", fmt.Errorf("chat request timeout")
		}
		return "", err
	}

	return result.Text(), err
}
