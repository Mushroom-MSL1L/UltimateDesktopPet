package chat

import (
	"context"
	"errors"

	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/database"
	pp "github.com/Mushroom-MSL1L/UltimateDesktopPet/pkg/print"

	"google.golang.org/genai"
)

const ChatCharLimit = 1000

type ChatMeta struct {
	Ctx             context.Context
	Controller      *database.BaseController[Dialog]
	DB              database.DB
	Dialog          *Dialog
	RolePlayContext string
	Client          *genai.Client
}

func init() {
	c := newChatController(nil)
	database.RegisterSchema(database.Pets, c)
	pp.Assert(pp.Chat, "chat init complete")
}

func newChatController(model **Dialog) *database.BaseController[Dialog] {
	return &database.BaseController[Dialog]{Model: model}
}

func NewChatMeta(ctx context.Context, geminiAPIKey string) *ChatMeta {
	var err error
	c := &ChatMeta{
		Ctx:    ctx,
		Dialog: &Dialog{},
	}
	c.Controller = newChatController(&c.Dialog)
	c.Client, err = newGeminiClient(geminiAPIKey)
	if err != nil {
		if errors.Is(err, errGeminiAPIKeyMissing) {
			pp.Warn(pp.Chat, "Gemini API key missing; set geminiAPIKey in System settings to enable chat")
		} else {
			pp.Fatal(pp.Chat, "newGeminiClient: err %v", err)
		}
	}
	return c
}

func (c *ChatMeta) Shutdown() {
	c.DB.CloseDB()
	pp.Assert(pp.Chat, "chat service stopped")
}
