package app

import (
	"cbot/pkg"
	"fmt"
	"strings"

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
		{
			Text:        "/search",
			Description: "search",
		},
		{
			Text:        "/create_course",
			Description: "create course",
		},
	}); err != nil {
		panic(err)
	}

	pkg.BOT.GetBot().Use(obj.validateUserMiddleware)

	pkg.BOT.GetBot().Handle("/start", obj.start)
	pkg.BOT.GetBot().Handle("/search", obj.searchCourseInWeb)
	pkg.BOT.GetBot().Handle("/create_course", obj.createCourseHandler)
	pkg.BOT.GetBot().Handle(telebot.OnText, obj.handleText)
	pkg.BOT.GetBot().Handle(telebot.OnCallback, obj.handleBtn)

	pkg.BOT.GetBot().Start()
}

func (obj *TGAppImpl) handleBtn(c telebot.Context) error {
	btnName := strings.TrimSpace(c.Callback().Data)

	switch btnName {

	}

	btnId := strings.Split(btnName, ":")

	if len(btnId) > 1 {
		switch btnId[0] {

		case "btn_course_name":
			return obj.setCourseName(c, btnId[1])

		case "btn_course_desc":
			return obj.setCourseDescription(c, btnId[1])

		case "btn_course_cost":
			return obj.setCourseCost(c, btnId[1])

		case "btn_course_duration":
			return obj.setCourseDuration(c, btnId[1])

		case "btn_course_link":
			return obj.setCourseLink(c, btnId[1])

		case "btn_course_unapprove":
			return obj.setCourseUnapprove(c, btnId[1])

		case "btn_course_approve":
			return obj.setCourseApprove(c, btnId[1])

		case "set_course_approve", "set_course_unapprove":
			return obj.handleCourseText(c, btnName)
		}
	}

	return c.Send("btn not found")
}

func (obj *TGAppImpl) handleText(c telebot.Context) error {
	var err error
	text := c.Message().Text
	cmd := pkg.CMDV.GetCommand(c.Sender().ID)

	switch cmd {
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

	default:
		if strings.Contains(cmd, "course") {
			err = obj.handleCourseText(c, cmd)
		}
	}

	if err != nil {
		return pkg.Trace(err)
	}

	pkg.CMDV.RemoveCommand(c.Sender().ID)

	user, err := pkg.USRV.Get(c.Sender().ID)
	if err != nil {
		return pkg.Trace(err)
	}

	return obj.validateUser(user)
}
