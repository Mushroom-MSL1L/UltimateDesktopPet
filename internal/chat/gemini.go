package chat

import (
	"context"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

const envFilePathForGemini = "configs/.env"
const model = "gemini-2.5-flash"

func newGeminiClient() (*genai.Client, error) {
	ctx := context.Background()
	err := godotenv.Load(envFilePathForGemini)
	if err != nil {
		return nil, err
	}
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}
	return client, err
}

func (c *ChatMeta) chatWithModel(ctx context.Context) (err error) {
	result, err := c.Client.Models.GenerateContent(
		ctx,
		model,
		genai.Text(c.Dialog.Request),
		nil,
	)
	if err != nil {
		c.Dialog.Response = ""
		return err
	}
	c.Dialog.Response = result.Text()
	return nil
}
