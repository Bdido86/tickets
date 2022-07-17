package commander

import (
	"fmt"
	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"ozon/go-hw-bot/config"
)

type Commander struct {
	bot *botApi.BotAPI
	//route map[string]CmdHandler
}

func Init(config *config.Config) (*Commander, error) {
	bot, err := botApi.NewBotAPI(config.Token())
	if err != nil {
		return nil, errors.Wrap(err, "init bot")
	}

	bot.Debug = config.Debug()
	fmt.Printf("Authorized on account %s", bot.Self.UserName)

	return &Commander{
		bot: bot,
		//route: make(map[string]CmdHandler),
	}, nil
}

func (c *Commander) Run() error {
	u := botApi.NewUpdate(0)
	u.Timeout = 30

	updates := c.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		cmd := update.Message.Command()

		fmt.Printf("[%s] %s <%s>", update.Message.From.UserName, update.Message.Text, cmd)

		msg := botApi.NewMessage(update.Message.Chat.ID, "test")
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := c.bot.Send(msg); err != nil {
			fmt.Print(err)
		}
	}

	return nil
}
