package app

import (
	"cbot/pkg"

	"gopkg.in/telebot.v4"
)

func (obj *TGAppImpl) start(c telebot.Context) error {
	sender := c.Sender()
	user := pkg.F.CreateUser()

	user.SetId(sender.ID)
	user.SetLogin(sender.Username)

	err := pkg.USRV.CreateOrUpdate(user)
	if err != nil{
		return pkg.Trace(err)
	}

	return nil
}
