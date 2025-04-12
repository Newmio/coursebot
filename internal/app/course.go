package app

import (
	"cbot/pkg"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/telebot.v4"
)

func (obj *TGAppImpl) searchCourseInWeb(c telebot.Context) error {
	parser := pkg.F.CreateCourseParser()

	courses, err := parser.StartParseSite("js", "prometheus")
	if err != nil {
		return c.Send(err.Error())
	}

	fmt.Println(courses)

	return nil
}

func (obj *TGAppImpl) createCourseHandler(c telebot.Context) error {
	return obj.createCourse(c, "")
}

func (obj *TGAppImpl) createCourse(c telebot.Context, id string) error {
	pkg.CMDV.SetCommand(c.Sender().ID, fmt.Sprintf("create_course_link:%s", id))
	return c.Send("Введіть посилання на курс")
}

func (obj *TGAppImpl) setCourseName(c telebot.Context, id string) error {
	pkg.CMDV.SetCommand(c.Sender().ID, fmt.Sprintf("set_course_name:%s", id))
	return c.Send("Введіть ім'я курсу")
}

func (obj *TGAppImpl) setCourseDescription(c telebot.Context, id string) error {
	pkg.CMDV.SetCommand(c.Sender().ID, fmt.Sprintf("set_course_desc:%s", id))
	return c.Send("Введіть опис курсу")
}

func (obj *TGAppImpl) setCourseCost(c telebot.Context, id string) error {
	pkg.CMDV.SetCommand(c.Sender().ID, fmt.Sprintf("set_course_cost:%s", id))
	return c.Send("Введіть ціну курсу")
}

func (obj *TGAppImpl) setCourseDuration(c telebot.Context, id string) error {
	pkg.CMDV.SetCommand(c.Sender().ID, fmt.Sprintf("set_course_duration:%s", id))
	return c.Send("Введіть тривалість курсу")
}

func (obj *TGAppImpl) setCourseLink(c telebot.Context, id string) error {
	return obj.createCourse(c, id)
}

func (obj *TGAppImpl) setCourseUnapprove(c telebot.Context, id string) error {
	inlineMenu := &telebot.ReplyMarkup{}
	inlineMenu.Inline([]telebot.Row{
		{
			telebot.Btn{
				Text:   "Так",
				Unique: fmt.Sprintf("set_course_unapprove:%s", id),
			},
			telebot.Btn{
				Text:   "Ні",
				Unique: fmt.Sprintf("set_course_approve:%s", id),
			},
		},
	}...)

	return c.Send("Зняти курс з підтвердження?", inlineMenu)
}

func (obj *TGAppImpl) setCourseApprove(c telebot.Context, id string) error {
	inlineMenu := &telebot.ReplyMarkup{}
	inlineMenu.Inline([]telebot.Row{
		{
			telebot.Btn{
				Text:   "Так",
				Unique: fmt.Sprintf("set_course_approve:%s", id),
			},
			telebot.Btn{
				Text:   "Ні",
				Unique: fmt.Sprintf("set_course_unapprove:%s", id),
			},
		},
	}...)

	return c.Send("Підтвердити курс?", inlineMenu)
}

func (obj *TGAppImpl) genCourseElems(c telebot.Context, link string, idCourse primitive.ObjectID) error {
	var course pkg.Course

	if !idCourse.IsZero() {
		if crv, ert := pkg.CRV.GetById(idCourse); ert != nil {
			return c.Send(pkg.Trace(ert).Error())
		} else {
			course = crv
		}
	} else if link != "" {
		if crv, ert := pkg.CRV.GetByLink(link); ert != nil {
			return c.Send(pkg.Trace(ert).Error())
		} else {
			course = crv
		}
	}

	id := course.GetId().Hex()

	apprBtn := telebot.Btn{}
	if course.GetApproved() {
		apprBtn.Text = "Зняти з підтвердження"
		apprBtn.Unique = fmt.Sprintf("btn_course_unapprove:%s", id)
	} else {
		apprBtn.Text = "Підтвердити"
		apprBtn.Unique = fmt.Sprintf("btn_course_approve:%s", id)
	}

	inlineMenu := &telebot.ReplyMarkup{}
	inlineMenu.Inline([]telebot.Row{
		{
			telebot.Btn{
				Text:   "Iм'я",
				Unique: fmt.Sprintf("btn_course_name:%s", id),
			},
			telebot.Btn{
				Text:   "Опис",
				Unique: fmt.Sprintf("btn_course_desc:%s", id),
			},
			telebot.Btn{
				Text:   "Ціна",
				Unique: fmt.Sprintf("btn_course_cost:%s", id),
			},
		},
		{
			telebot.Btn{
				Text:   "Тривалість",
				Unique: fmt.Sprintf("btn_course_duration:%s", id),
			},
			telebot.Btn{
				Text:   "Посилання",
				Unique: fmt.Sprintf("btn_course_link:%s", id),
			},
		},
		{
			apprBtn,
		},
	}...)

	if course.GetName() == "" {
		course.SetName("-")
	}

	if course.GetDescription() == "" {
		course.SetDescription("-")
	}

	if course.GetDuration() == "" {
		course.SetDuration("-")
	}

	if course.GetLink() == "" {
		course.SetLink("-")
	}

	text := fmt.Sprintf("Назва: %s\nОпис: %s\nЦіна: %s\nТривалість: %s\nПосилання: %s",
		course.GetName(), course.GetDescription(), course.GetCost(), course.GetDuration(), course.GetLink())

	return c.Send(text, inlineMenu)
}

func (obj *TGAppImpl) handleCourseText(c telebot.Context, cmd string) error {
	var err error
	var idStr string
	var id primitive.ObjectID
	var needUpdate, needShow bool

	text := c.Message().Text
	course := pkg.F.CreateCourse()

	parts := strings.Split(cmd, ":")
	if len(parts) > 1 {
		idStr = parts[1]
	}
	cmd = parts[0]

	if idStr != "" {
		if objId, ert := primitive.ObjectIDFromHex(idStr); ert != nil {
			return pkg.Trace(ert)
		} else {
			id = objId
		}
	}

	if !id.IsZero() {
		if crv, ert := pkg.CRV.GetById(id); ert != nil {
			return pkg.Trace(ert)
		} else {
			course = crv
		}
	}

	switch cmd {
	case "create_course_link":
		if exists, ert := pkg.CRV.Exists(text); ert != nil {
			err = ert
		} else {
			needShow = true

			if !exists {
				course.SetLink(text)
				needUpdate = true
			}
		}

	case "set_course_name":
		course.SetName(text)
		needUpdate = true
		needShow = true

	case "set_course_desc":
		course.SetDescription(text)
		needUpdate = true
		needShow = true

	case "set_course_cost":
		course.SetCost(text)
		needUpdate = true
		needShow = true

	case "set_course_duration":
		needUpdate = true
		needShow = true

	case "set_course_approve":
		course.SetApproved(true)
		needUpdate = true
		needShow = true

	case "set_course_unapprove":
		course.SetApproved(false)
		needUpdate = true
		needShow = true
	}

	if needUpdate {
		if ert := pkg.CRV.CreateOrUpdate(course); ert != nil {
			err = ert
		}
	}

	if needShow && err == nil {
		err = obj.genCourseElems(c, course.GetLink(), course.GetId())
	}

	if err != nil {
		return pkg.Trace(err)
	}

	return nil
}
