package pkg

import "gopkg.in/telebot.v4"

type Bot interface {
	GetBot() *telebot.Bot
}
