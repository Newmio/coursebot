package app

import (
	"cbot/pkg"
	"fmt"
	"strings"

	"gopkg.in/telebot.v4"
)

func (obj *TGAppImpl) validateUserMiddleware(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {

		user, err := pkg.USRV.Get(c.Sender().ID)
		if err != nil {
			return obj.start(c)
		} else {
			c.Set("user", user)
		}

		cmd := pkg.CMDV.GetCommand(c.Sender().ID)

		if user != nil &&
			!strings.Contains(cmd, "set_first_name") &&
			!strings.Contains(cmd, "set_last_name") &&
			!strings.Contains(cmd, "set_middle_name") {
			if err := obj.validateUser(user); err != nil {
				return pkg.BOT.Send(c, true, err.Error())
			}
		}

		needClear := true

		if c.Callback() != nil {
			msg := c.Callback().Message
			markup := msg.ReplyMarkup

			if markup != nil {

				for _, row := range markup.InlineKeyboard {
					for _, btn := range row {
						if strings.Contains(btn.Data, "btn_clear_msg") {
							needClear = false
						}
					}
				}
			}
		}

		if needClear {
			pkg.CMDV.ClearAllDeleteMessages(c.Sender().ID)
		}

		if err := next(c); err == nil {
			user := c.Get("user").(pkg.User)
			
			if err := obj.validateUser(user); err != nil {
				return pkg.BOT.Send(c, true, err.Error())
			}
			return nil
		}else{
			return err
		}
	}
}

func (obj *TGAppImpl) validateUser(user pkg.User) error {
	if user.GetFirstName() == "" {
		pkg.CMDV.SetCommand(user.GetId(), "set_first_name")
		return fmt.Errorf("*Ваше ім'я не вказано. Введіть своє ім'я")
	}

	if user.GetLastName() == "" {
		pkg.CMDV.SetCommand(user.GetId(), "set_last_name")
		return fmt.Errorf("*Ваше прізвище не вказано. Введіть своє прізвище")
	}

	if user.GetMiddleName() == "" {
		pkg.CMDV.SetCommand(user.GetId(), "set_middle_name")
		return fmt.Errorf("*Ваше по-батькове не вказано. Введіть своє по-батькове")
	}

	return nil
}
