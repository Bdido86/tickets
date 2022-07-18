package commander

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type Message struct {
	cmd       string
	userName  string
	userId    uint
	arguments []string
}

func newMessage(m *tgbotapi.Message) Message {
	parsedCmd := strings.Split(m.Command(), "_")
	message := Message{
		cmd:      parsedCmd[0],
		userName: m.From.UserName,
		userId:   uint(m.From.ID),
	}

	if len(parsedCmd) > 1 {
		message.arguments = parsedCmd[1:]
	}
	return message
}

func (m Message) Cmd() string {
	return m.cmd
}

func (m Message) Arguments() []string {
	return m.arguments
}

func (m Message) UserName() string {
	return m.userName
}

func (m Message) UserId() uint {
	return m.userId
}
