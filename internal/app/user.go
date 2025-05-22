package app

import (
	"cbot/pkg"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/telebot.v4"
)

func (obj *TGAppImpl) setCourseCoins(c telebot.Context, courseIdStr, userIdStr string) error {
	text := c.Message().Text

	coin, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		return pkg.BOT.Send(c, false, err.Error())
	}

	courseId, err := primitive.ObjectIDFromHex(courseIdStr)
	if err != nil {
		return pkg.BOT.Send(c, false, err.Error())
	}

	userId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		return pkg.BOT.Send(c, false, err.Error())
	}

	if err := pkg.CRV.SetCoins(courseId, userId, coin); err != nil {
		return pkg.BOT.Send(c, false, err.Error())
	}

	coins, err := pkg.CRV.GetCoins(courseId, userId)
	if err != nil {
		return pkg.BOT.Send(c, false, err.Error())
	}

	admins, err := pkg.USRV.GetAdmins()
	if err != nil {
		return pkg.BOT.Send(c, false, err.Error())
	}

	var resultCoin int

	if len(coins) == len(admins) {
		var sum int

		for _, coin := range coins {
			sum += coin
		}

		resultCoin = sum / len(coins)

		if err := pkg.CRV.SetResultCoins(courseId, userId, resultCoin); err != nil {
			return pkg.BOT.Send(c, false, err.Error())
		}

		if err := pkg.CRV.StopCourse(courseId, userId); err != nil {
			return pkg.BOT.Send(c, false, err.Error())
		}

		if err := pkg.CRV.UpdateCheckAdmin(courseId, userId, false); err != nil {
			return pkg.BOT.Send(c, false, err.Error())
		}
	}

	return pkg.BOT.Send(c, false, fmt.Sprintf("Кiнцевиий бал (%d) успішно встановлено", resultCoin))
}

func (obj *TGAppImpl) sendCourseCoins(c telebot.Context, courseIdStr, userIdStr string) error {
	pkg.CMDV.SetCommand(c.Sender().ID, fmt.Sprintf("set_sendcoin:%s:%s", courseIdStr, userIdStr))
	return pkg.BOT.Send(c, false, "Введіть кiлькiсть балiв (ціле число)\nКоли всi балi бiли виставленi\nСтуденту буде виставлено середнє значення\n")
}

func (obj *TGAppImpl) start(c telebot.Context) error {
	sender := c.Sender()
	user := pkg.F.CreateUser()

	var dbUser pkg.User

	if dU, ok := c.Get("user").(pkg.User); ok {
		dbUser = dU
	}

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

	user, err = pkg.USRV.Get(sender.ID)
	if err != nil {
		return pkg.Trace(err)
	}

	c.Set("user", user)

	return obj.menu(c)
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

	if user.GetIsAdmin() {
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
