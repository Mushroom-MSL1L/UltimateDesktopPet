package chat

import (
	pp "UltimateDesktopPet/pkg/print"
	"strings"
)

func (c *ChatMeta) ChatWithPet(userInput string) (string, error) {
	request := strings.TrimSpace(userInput)
	if len(request) > ChatCharLimit {
		request = request[:ChatCharLimit]
		pp.Warn(pp.Chat, "Chat request too long, truncated with %d chars", ChatCharLimit)
	}

	pp.Info(pp.App, "ChatWithPet: received message %q", request)
	response, err := c.Chat(request)
	if err != nil {
		pp.Warn(pp.Chat, "Chat request \"%s\" with error: %v", request, err)
		return "", err
	}
	return response, nil
}
