package app

import (
	"cbot/pkg"
	"strings"

	"gopkg.in/telebot.v4"
)

func (obj *TGAppImpl) ValidateUserMiddleware(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		user, err := pkg.USRV.Get(c.Sender().ID)
		if err != nil {
			return pkg.Trace(err)
		}

		text := c.Message().Text

		if user != nil && strings.HasPrefix(text, "/") {
			if user.GetFirstName() == "" {
				pkg.CMDV.SetCommand(user.GetId(), "set_first_name")
				return c.Send("Ваше ім'я не вказано. Введіть своє ім'я")
			}

			if user.GetLastName() == "" {
				pkg.CMDV.SetCommand(user.GetId(), "set_last_name")
				return c.Send("Ваше прізвище не вказано. Введіть своє прізвище")
			}

			if user.GetMiddleName() == "" {
				pkg.CMDV.SetCommand(user.GetId(), "set_middle_name")
				return c.Send("Ваше по-батькове не вказано. Введіть своє по-батькове")
			}
		}

		c.Set("user", user)
		return next(c)
	}
}
