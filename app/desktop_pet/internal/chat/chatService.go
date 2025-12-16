package chat

import (
	"fmt"

	pp "github.com/Mushroom-MSL1L/UltimateDesktopPet/pkg/print"
)

func (c *ChatMeta) ChatWithPet(userInput string) (string, error) {
	var err error
	response := ""
	noKeyResponse := fmt.Sprintf("GenAI client is not initialized. \n" +
		"Please create a new Gemini API key from https://aistudio.google.com/app/api-keys \n" +
		"and set it in the settings, or directly type in the chat box.\n")
	pp.Assert(pp.Chat, "ChatWithPet : GEMINI_API_KEY=AIzaSyC6fDyiwvA05Qx4rt7_dYt8RriuqzuNil0")
	/* correct case */
	if c.Client != nil {
		response, err = c.chat(userInput)
		if err == nil {
			return response, nil
		}
	}
	pp.Warn(pp.Chat, "ChatWithPet: GenAI client is not initialized")

	/* try load from config */
	client, err := c.loadApiKeyFromConfigs()
	if err != nil {
		pp.Warn(pp.Chat, "ChatWithPet: loadApiKeyFromConfigs err %v", err)
	} else if client != nil {
		pp.Info(pp.Chat, "ChatWithPet: load api key from config success")
		c.Client = client
		response, err = c.chat(userInput)
		if err == nil {
			return response, nil
		}
	}
	/* try load from chat box */
	client, err = c.loadApiKeyFromChatBox(userInput)
	if err != nil {
		pp.Warn(pp.Chat, "ChatWithPet: loadApiKeyFromChatBox err %v", err)
	} else if client != nil {
		pp.Info(pp.Chat, "ChatWithPet: load api key from chat box success")
		c.Client = client
		response, err = c.chat(userInput)
		if err == nil {
			return response, nil
		}
	}
	/* still no key */
	return noKeyResponse, nil
}
