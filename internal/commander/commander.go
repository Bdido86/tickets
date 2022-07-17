package commander

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"ozon/go-hw-bot/config"
)

type CmdHandler func(m Message) string

func (c *Commander) RegisterHandler(f CmdHandler) {
	c.route = f
}

type Commander struct {
	bot   *tgbotapi.BotAPI
	route CmdHandler
}

func Init(config *config.Config) (*Commander, error) {
	bot, err := tgbotapi.NewBotAPI(config.Token())
	if err != nil {
		return nil, errors.Wrap(err, "init bot")
	}

	bot.Debug = config.Debug()
	fmt.Printf("Authorized on account %s", bot.Self.UserName)

	return &Commander{
		bot: bot,
	}, nil
}

func (c *Commander) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

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
		msg.Text = c.route(message)

		if _, err := c.bot.Send(msg); err != nil {
			fmt.Print(err)
		}
	}

	return nil
}
