package commander

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type Message struct {
	cmd       string
	arguments []string
	userName  string
}

func newMessage(m *tgbotapi.Message) Message {
	message := Message{
		cmd:      m.Command(),
		userName: m.From.UserName,
	}

	arg := m.CommandArguments()
	if arg != "" {
		message.arguments = strings.Split(arg, " ")
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

//func (m Message) CallBackData() string {
//	return m.callBackData
//}
//
//func (m Message) IsCallBackData() bool {
//	return m.callBackData != "" && len(m.callBackData) > 0
//}
