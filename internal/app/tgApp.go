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
			Description: "Старт",
		},
		{
			Text:        "/menu",
			Description: "Меню",
		},
		{
			Text:        "/profile",
			Description: "Профіль",
		},
		{
			Text:        "/my_courses",
			Description: "Мої курси",
		},
		{
			Text:        "/completed",
			Description: "Завершені курси",
		},
		{
			Text:        "/search",
			Description: "Пошук",
		},
		{
			Text:        "/create_course",
			Description: "Створити курс",
		},
		{
			Text:        "/search_green_courses",
			Description: "Знайти підтверджені курси",
		},
	}); err != nil {
		panic(err)
	}

	pkg.BOT.GetBot().Use(obj.validateUserMiddleware)

	pkg.BOT.GetBot().Handle("/start", obj.start)
	pkg.BOT.GetBot().Handle("/search", obj.searchCourseInWeb)
	pkg.BOT.GetBot().Handle("/create_course", obj.createCourseHandler)
	pkg.BOT.GetBot().Handle("/search_green_courses", obj.getApprovedCourses)
	pkg.BOT.GetBot().Handle("/profile", obj.profile)
	pkg.BOT.GetBot().Handle("/my_courses", obj.myCourses)
	pkg.BOT.GetBot().Handle("/menu", obj.menu)
	pkg.BOT.GetBot().Handle("/completed", obj.getCompletedCourses)
	pkg.BOT.GetBot().Handle(telebot.OnText, obj.handleText)
	pkg.BOT.GetBot().Handle(telebot.OnCallback, obj.handleBtn)
	pkg.BOT.GetBot().Handle(telebot.OnDocument, obj.handleText)
	pkg.BOT.GetBot().Handle(telebot.OnPhoto, obj.handleText)

	pkg.BOT.GetBot().Start()
}

func (obj *TGAppImpl) handleBtn(c telebot.Context) error {
	btnName := strings.TrimSpace(c.Callback().Data)

	switch btnName {
	case "btn_profile":
		return obj.profile(c)

	case "btn_my_courses":
		return obj.myCourses(c)

	case "btn_completed":
		return obj.getCompletedCourses(c)

	case "btn_search_green_courses":
		return obj.getApprovedCourses(c)

	case "btn_create_course":
		return obj.createCourseHandler(c)

	case "btn_search_red_courses":
		//TODO

	case "btn_clear_msg":
		msg := c.Callback().Message

		return pkg.BOT.GetBot().Delete(msg)
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

		case "btn_start_course":
			return obj.startCourse(c, btnId[1])

		case "btn_stop_course":
			return obj.stopCourse(c, btnId[1])

		case "set_course_approve", "set_course_unapprove":
			return obj.handleCourseText(c, btnName)

		case "get_approved_courses":
			return obj.handleCourseText(c, btnName)

		case "btn_send_result_course":
			return obj.sendCourseResult(c, btnId[1])

		case "btn_d_c_f":
			return obj.deleteFilesByCourse(c, btnId[1], btnId[2])

		case "btn_complete_send_result_course":
			return obj.sendCheckResultToAdmin(c, btnId[1])

		case "btn_sendcoin":
			return obj.sendCourseCoins(c, btnId[1], btnId[2])

		case "btn_c_f":
			return obj.getCourseFiles(c, btnId[1], btnId[2])

		case "btn_get_course":
			return obj.getCourse(c, btnId[1])

		case "btn_get_course_desc":
			return obj.getCourseDescription(c, btnId[1])
		}
	}

	return pkg.BOT.Send(c, false, "btn not found")
}

func (obj *TGAppImpl) handleText(c telebot.Context) error {
	var err error
	text := c.Message().Text
	cmd := pkg.CMDV.GetCommand(c.Sender().ID)
	pkg.CMDV.RemoveCommand(c.Sender().ID)

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
		if strings.Contains(cmd, "course") ||
			strings.Contains(cmd, "set_sendcoin") {
			err = obj.handleCourseText(c, cmd)
		}
	}

	if err != nil {
		return pkg.Trace(err)
	}

	user, err := pkg.USRV.Get(c.Sender().ID)
	if err != nil {
		return pkg.Trace(err)
	}

	c.Set("user", user)

	return nil
}

func (obj *TGAppImpl) menu(c telebot.Context) error {
	user := c.Get("user").(pkg.User)
	inlineMenu := &telebot.ReplyMarkup{}

	if user.GetIsAdmin() {
		inlineMenu.Inline([]telebot.Row{
			{
				telebot.Btn{
					Text:   "🖥️Профіль🖥️",
					Unique: "btn_profile",
				},
			},
			{
				telebot.Btn{
					Text:   "🔨Створити курс🔨",
					Unique: "btn_create_course",
				},
			},
			{
				telebot.Btn{
					Text:   "🟢Знайти підтверджені курси🟢",
					Unique: "btn_search_green_courses",
				},
			},
			{
				telebot.Btn{
					Text:   "🔴Знайти не підтверджені курси🔴",
					Unique: "btn_search_red_courses",
				},
			},
		}...)

	} else {
		inlineMenu.Inline([]telebot.Row{
			{
				telebot.Btn{
					Text:   "🖥️Профіль🖥️",
					Unique: "btn_profile",
				},
				telebot.Btn{
					Text:   "🔄Мої курси🔄",
					Unique: "btn_my_courses",
				},
			},
			{
				telebot.Btn{
					Text:   "✔️Завершені курси✔️",
					Unique: "btn_completed",
				},
			},
			{
				telebot.Btn{
					Text:   "🟢Знайти підтверджені курси🟢",
					Unique: "btn_search_green_courses",
				},
			},
			{
				telebot.Btn{
					Text:   "🔴Знайти не підтверджені курси🔴",
					Unique: "btn_search_red_courses",
				},
			},
		}...)
	}

	return pkg.BOT.Send(c, true, "💎Меню💎", inlineMenu)
}
