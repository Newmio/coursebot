package core

import (
	"cbot/pkg"

	"gopkg.in/telebot.v4"
)

type BotImpl struct {
	mBot *telebot.Bot
}

func CreateBot() pkg.Bot {
	b, err := telebot.NewBot(telebot.Settings{
		Token:  pkg.BotToken,
		Poller: &telebot.LongPoller{Timeout: 2},
	})
	if err != nil {
		panic(err)
	}

	return &BotImpl{mBot: b}
}

func (obj *BotImpl) GetBot() *telebot.Bot {
	return obj.mBot
}