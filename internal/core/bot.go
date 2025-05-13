package core

import (
	"cbot/pkg"

	"gopkg.in/telebot.v4"
)

type BotImpl struct {
	mBot *telebot.Bot

	mReqCount int //TODO сделать лоад балансер
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

func (obj *BotImpl) Send(c telebot.Context, needDel bool, what interface{}, opts ...interface{}) error {
	var err error

	if needDel {
		if msg, ert := obj.mBot.Send(c.Sender(), what, opts...); ert == nil {
			pkg.CMDV.AppendDeleteMessage(msg)
		}else{
			err = ert
		}
	} else {
		err = c.Send(what, opts...)
	}

	return err
}
