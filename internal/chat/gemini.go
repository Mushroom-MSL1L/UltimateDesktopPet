package chat

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

const envFilePathForGemini = "configs/.env"
const model = "gemini-2.5-flash"

func newGeminiClient() (*genai.Client, error) {
	ctx := context.Background()
	err := godotenv.Load(envFilePathForGemini)
	if err != nil {
		err = fmt.Errorf(
			"newGeminiClient: \n\tfailed to load .env file: %w \n\n\t"+
				"We need user to provide a valid genimi api at %s, content like \"GEMINI_API_KEY=<your_api_key>\"\n\t"+
				"You can generate a key by the link https://aistudio.google.com/app/api-keys",
			err, envFilePathForGemini)
		return nil, err
	}
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		err = fmt.Errorf(
			"newGeminiClient: \n\tfailed to load .env file: %w \n\n\t"+
				"We need user to provide a valid genimi api at %s, content like \"GEMINI_API_KEY=<your_api_key>\"\n\t"+
				"You can generate a key by the link https://aistudio.google.com/app/api-keys",
			err, envFilePathForGemini)
		return nil, err
	}
	return client, err
}

func (c *ChatMeta) chatWithModel(ctx context.Context, request string) (response string, err error) {
	result, err := c.Client.Models.GenerateContent(
		ctx,
		model,
		genai.Text(request),
		nil,
	)
	return result.Text(), err
}
