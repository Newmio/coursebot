package app

import (
	"cbot/pkg"
	"fmt"

	"gopkg.in/telebot.v4"
)

type TGAppImpl struct{}

func CreateTGApp() pkg.TGApp {
	return &TGAppImpl{}
}

func (obj *TGAppImpl) Run() {
	if err := pkg.BOT.GetBot().SetCommands([]telebot.Command{
		{
			Text:        "/start",
			Description: "start",
		},
	}); err != nil {
		panic(err)
	}

	pkg.BOT.GetBot().Use(obj.ValidateUserMiddleware)

	pkg.BOT.GetBot().Handle("/start", obj.start)
	pkg.BOT.GetBot().Handle(telebot.OnText, obj.HandleText)

	pkg.BOT.GetBot().Start()
}

func (obj *TGAppImpl) HandleText(c telebot.Context) error {
	var err error
	text := c.Message().Text

	switch pkg.CMDV.GetCommand(c.Sender().ID) {
	case "set_first_name":
		if user, ok := c.Get("user").(pkg.User); ok {
			user.SetFirstName(text)
			err = pkg.USRV.CreateOrUpdate(user)
		} else {
			err = fmt.Errorf("error cast user")
		}

	case "set_middle_name":
		if user, ok := c.Get("user").(pkg.User); ok {
			user.SetMiddleName(text)
			err = pkg.USRV.CreateOrUpdate(user)
		} else {
			err = fmt.Errorf("error cast user")
		}

	case "set_last_name":
		if user, ok := c.Get("user").(pkg.User); ok {
			user.SetLastName(text)
			err = pkg.USRV.CreateOrUpdate(user)
		} else {
			err = fmt.Errorf("error cast user")
		}
	}

	if err != nil {
		return pkg.Trace(err)
	}

	pkg.CMDV.RemoveCommand(c.Sender().ID)

	return nil
}
