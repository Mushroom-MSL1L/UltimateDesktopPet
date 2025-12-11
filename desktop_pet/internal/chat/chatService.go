package chat

import (
	"strings"
	"time"

	pp "github.com/Mushroom-MSL1L/UltimateDesktopPet/desktop_pet/pkg/print"
)

func (c *ChatMeta) ChatWithPet(userInput string) (string, error) {
	request := strings.TrimSpace(userInput)
	if len(request) > ChatCharLimit {
		request = request[:ChatCharLimit]
		pp.Warn(pp.Chat, "Chat request too long, truncated with %d chars", ChatCharLimit)
	}
	pp.Info(pp.Chat, "Chat: received message %q", request)

	c.Dialog.Request = request
	c.Dialog.Timestamp = time.Now()
	response, err := c.chatWithModel(c.Ctx, petRolePlayPrompt(c.RolePlayContext, request))
	if err != nil {
		pp.Warn(pp.Chat, "Chat request \"%s\" with error: %v", request, err)
		return "", err
	}
	c.Dialog.Response = response
	c.Controller.Create(c.DB.GetDB())

	return response, nil
}
