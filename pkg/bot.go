package pkg

import "gopkg.in/telebot.v4"

type Bot interface {
	GetBot() *telebot.Bot
	Send(c telebot.Context, needDel bool, what interface{}, opts ...interface{}) error
}
