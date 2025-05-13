package pkg

import "gopkg.in/telebot.v4"

type CommandVault interface {
	SetCommand(userId int64, command string)
	GetCommand(userId int64) string
	RemoveCommand(userId int64)
	AppendDeleteMessage(msg *telebot.Message)
	ClearAllDeleteMessages(userId int64)
}
