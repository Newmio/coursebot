package app

import (
	"cbot/pkg"

	"gopkg.in/telebot.v4"
)

func (obj *TGAppImpl) start(c telebot.Context) error {
	sender := c.Sender()
	user := pkg.F.CreateUser()

	dbUser := c.Get("user").(pkg.User)
	if dbUser != nil {
		user.SetIsAdmin(dbUser.GetIsAdmin())
	} else {
		user.SetIsAdmin(false)
	}

	user.SetId(sender.ID)
	user.SetLogin(sender.Username)

	err := pkg.USRV.CreateOrUpdate(user)
	if err != nil {
		return pkg.Trace(err)
	}

	return nil
}

func (obj *TGAppImpl) profile(c telebot.Context) error {
	user := c.Get("user").(pkg.User)

	inlineMenu := &telebot.ReplyMarkup{}
	rows := []telebot.Row{
		{
			telebot.Btn{
				Text:   "Iм'я",
				Unique: "set_first_name",
			},
			telebot.Btn{
				Text:   "Прізвище",
				Unique: "set_last_name",
			},
			telebot.Btn{
				Text:   "По-батькові",
				Unique: "set_middle_name",
			},
		},
	}

	if user.GetIsAdmin(){
		rows = append(rows, []telebot.Row{
			{
				telebot.Btn{
					Text:   "Група",
					Unique: "btn_set_group",
				},
				telebot.Btn{
					Text:   "Група",
					Unique: "btn_set_group_num",
				},
			},
			{
				telebot.Btn{
					Text:   "Професія",
					Unique: "btn_set_proffesion",
				},
				telebot.Btn{
					Text:   "Номер професіі",
					Unique: "btn_set_proffesion_num",
				},
			},
		}...)
	}
	inlineMenu.Inline(rows...)

	return pkg.BOT.Send(c, true, user.String(), inlineMenu)
}
