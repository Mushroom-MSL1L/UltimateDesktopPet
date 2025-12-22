package chat

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

var errGeminiAPIKeyMissing = errors.New("gemini API key missing")

const model = "gemini-2.5-flash"

func newGeminiClient(apiKey string) (*genai.Client, error) {
	key := strings.TrimSpace(apiKey)
	if key == "" {
		return nil, fmt.Errorf("%w: set geminiAPIKey in system.yaml from the System settings dialog", errGeminiAPIKeyMissing)
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: key})
	if err != nil {
		return nil, fmt.Errorf("newGeminiClient: failed to init Gemini client: %w", err)
	}
	return client, err
}

func (c *ChatMeta) chatWithModel(ctx context.Context, request string) (response string, err error) {
	if c.Client == nil {
		return "", fmt.Errorf("%w: open System settings and add your key to geminiAPIKey", errGeminiAPIKeyMissing)
	}
	result, err := c.Client.Models.GenerateContent(
		ctx,
		model,
		genai.Text(request),
		nil,
	)
	return result.Text(), err
}
