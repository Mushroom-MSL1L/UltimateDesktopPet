package chat

import (
	"context"
	"strings"
	"time"

	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/configs"
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/desktop_pet/internal/database"
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/pkg/configLoader"
	pp "github.com/Mushroom-MSL1L/UltimateDesktopPet/pkg/print"

	"google.golang.org/genai"
)

const ChatCharLimit = 1000
const envFilePathForGemini = "./configs/.env"
const model = "gemini-2.5-flash"
const geminiAPIConfigKey = "geminiApiKey"

type ChatMeta struct {
	Ctx             context.Context
	Controller      *database.BaseController[Dialog]
	DB              database.DB
	Dialog          *Dialog
	RolePlayContext string
	Client          *genai.Client
	configPath      string
}

func init() {
	c := newChatController(nil)
	database.RegisterSchema(database.Pets, c)
	pp.Assert(pp.Chat, "chat init complete")
}

func newChatController(model **Dialog) *database.BaseController[Dialog] {
	return &database.BaseController[Dialog]{Model: model}
}

func NewChatMeta(ctx context.Context, configPath string, apiKey string, rolePlayContext string) *ChatMeta {
	var err error
	c := &ChatMeta{
		Ctx:    ctx,
		Dialog: &Dialog{},
	}
	c.Controller = newChatController(&c.Dialog)
	c.configPath = configPath
	c.RolePlayContext = rolePlayContext

	c.Client, err = newGeminiClient(apiKey)
	if err != nil {
		pp.Warn(pp.Chat, "newGeminiClient: err %v", err)
	}
	pp.Warn(pp.Chat, "newChatMeta: Gemini client initialized %v", c.Client != nil)
	return c
}

func (c *ChatMeta) Shutdown() {
	c.DB.CloseDB()
	pp.Assert(pp.Chat, "chat service stopped")
}

func (c *ChatMeta) chat(userInput string) (string, error) {
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

func (c *ChatMeta) loadApiKeyFromConfigs() (*genai.Client, error) {
	var err error
	cfgs := configLoader.LoadConfig(c.configPath, &configs.System{})

	c.Client, err = newGeminiClient(cfgs.GeminiAPIKey)
	if err != nil {
		pp.Warn(pp.Chat, "loadApiKeyFromConfigs: err %v", err)
	}
	return c.Client, nil
}

func (c *ChatMeta) loadApiKeyFromChatBox(userInput string) (*genai.Client, error) {
	var err error
	configLoader.UpdateConfig(c.configPath, &configs.System{}, geminiAPIConfigKey, userInput)
	cfgs := configLoader.LoadConfig(c.configPath, &configs.System{})

	c.Client, err = newGeminiClient(cfgs.GeminiAPIKey)
	if err != nil {
		pp.Warn(pp.Chat, "loadApiKeyFromChatBox: err %v", err)
	}
	return c.Client, nil
}
