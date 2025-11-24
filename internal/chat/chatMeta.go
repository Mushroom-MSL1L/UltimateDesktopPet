package chat

import (
	"UltimateDesktopPet/internal/database"
	pp "UltimateDesktopPet/pkg/print"
	"context"

	"google.golang.org/genai"
)

const ChatCharLimit = 1000

type ChatMeta struct {
	Context    context.Context
	Controller *database.BaseController[Dialog]
	DB         database.DB
	Dialog     *Dialog
	Client     *genai.Client
}

func init() {
	c := newChatController(nil)
	database.RegisterSchema(database.Pets, c)
	pp.Assert(pp.Chat, "chat init complete")
}

func newChatController(model **Dialog) *database.BaseController[Dialog] {
	return &database.BaseController[Dialog]{Model: model}
}

func NewChatMeta(ctx context.Context) *ChatMeta {
	var err error
	c := &ChatMeta{
		Context: ctx,
		Dialog:  &Dialog{},
	}
	c.Controller = newChatController(&c.Dialog)
	c.Client, err = newGeminiClient()
	if err != nil {
		pp.Fatal(pp.Chat, "newGeminiClient: err %v", err)
	}
	return c
}

func (c *ChatMeta) Shutdown() {
	c.DB.CloseDB()
	pp.Assert(pp.Chat, "chat service stopped")
}

func (c *ChatMeta) Chat(request string) (response string, err error) {
	c.Dialog.Request = request
	err = c.chatWithModel(c.Context)
	if err != nil {
		return "", err
	}
	c.Controller.Create(c.DB.GetDB())
	return c.Dialog.Response, nil
}
