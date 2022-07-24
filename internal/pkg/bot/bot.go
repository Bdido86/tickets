package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/config"
)

const botTimeout int = 60

type CommandProcessor func(m Message) string

func (c *Commander) RegisterCommandProcessor(cp CommandProcessor) {
	c.processor = cp
}

type Commander struct {
	bot       *tgbotapi.BotAPI
	processor CommandProcessor
}

func Init(config *config.Config) (*Commander, error) {
	bot, err := tgbotapi.NewBotAPI(config.Token())
	if err != nil {
		return nil, errors.Wrap(err, "init bot")
	}

	bot.Debug = config.Debug()

	return &Commander{
		bot: bot,
	}, nil
}

func (c *Commander) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = botTimeout

	updates := c.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		message := newMessage(update.Message)
		msg.Text = c.processor(message)

		if _, err := c.bot.Send(msg); err != nil {
			fmt.Print(err)
		}
	}

	return nil
}