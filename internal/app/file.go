package app

import (
	"cbot/pkg"

	"gopkg.in/telebot.v4"
)

func (obj *TGAppImpl) courseResultFileHandler(c telebot.Context) error {
	msg := c.Message()
	
	if msg.Document != nil{
		
	}else if msg.Photo != nil{

	}else{
		return pkg.BOT.Send(c, true, "Невідомий тип повідомлення")
	}

	return nil
}
