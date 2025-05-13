package app

import (
	"cbot/pkg"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/telebot.v4"
)

func (obj *TGAppImpl) getApprovedCourses(c telebot.Context) error {
	pkg.CMDV.SetCommand(c.Sender().ID, "get_approved_courses:0")
	return pkg.BOT.Send(c, true, "Введіть назву курсу для пошуку")
}

func (obj *TGAppImpl) searchCourseInWeb(c telebot.Context) error {
	parser := pkg.F.CreateCourseParser()

	courses, err := parser.StartParseSite("js", "prometheus")
	if err != nil {
		return pkg.BOT.Send(c, false, err.Error())
	}

	fmt.Println(courses)

	return nil
}

func (obj *TGAppImpl) createCourseHandler(c telebot.Context) error {
	return obj.createCourse(c, "")
}

func (obj *TGAppImpl) createCourse(c telebot.Context, id string) error {
	pkg.CMDV.SetCommand(c.Sender().ID, fmt.Sprintf("create_course_link:%s", id))
	return pkg.BOT.Send(c, true, "Введіть посилання на курс")
}

func (obj *TGAppImpl) setCourseName(c telebot.Context, id string) error {
	pkg.CMDV.SetCommand(c.Sender().ID, fmt.Sprintf("set_course_name:%s", id))
	return pkg.BOT.Send(c, true, "Введіть ім'я курсу")
}

func (obj *TGAppImpl) setCourseDescription(c telebot.Context, id string) error {
	pkg.CMDV.SetCommand(c.Sender().ID, fmt.Sprintf("set_course_desc:%s", id))
	return pkg.BOT.Send(c, true, "Введіть опис курсу")
}

func (obj *TGAppImpl) setCourseCost(c telebot.Context, id string) error {
	pkg.CMDV.SetCommand(c.Sender().ID, fmt.Sprintf("set_course_cost:%s", id))
	return pkg.BOT.Send(c, true, "Введіть ціну курсу")
}

func (obj *TGAppImpl) setCourseDuration(c telebot.Context, id string) error {
	pkg.CMDV.SetCommand(c.Sender().ID, fmt.Sprintf("set_course_duration:%s", id))
	return pkg.BOT.Send(c, true, "Введіть тривалість курсу")
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

	return pkg.BOT.Send(c, true, "Зняти курс з підтвердження?", inlineMenu)
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

	return pkg.BOT.Send(c, true, "Підтвердити курс?", inlineMenu)
}

func getTunerBtns(course pkg.Course) *telebot.ReplyMarkup {
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

	return inlineMenu
}

func (obj *TGAppImpl) handleCourseText(c telebot.Context, cmd string) error {
	var err error
	var paramStr string
	var id primitive.ObjectID
	var needUpdate, needShow, showMany bool

	pkg.CMDV.ClearAllDeleteMessages(c.Sender().ID)

	text := c.Message().Text
	course := pkg.F.CreateCourse()

	parts := strings.Split(cmd, ":")
	if len(parts) > 1 {
		paramStr = parts[1]
	}
	cmd = parts[0]

	if paramStr != "" {
		if objId, ert := primitive.ObjectIDFromHex(paramStr); ert == nil {
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
		course.SetDuration(text)
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

	case "get_approved_courses":
		course.SetApproved(true)
		showMany = true
	}

	if needUpdate {
		if ert := pkg.CRV.CreateOrUpdate(course); ert != nil {
			err = ert
		}
	}

	if needShow && err == nil {
		link := course.GetLink()
		idCourse := course.GetId()

		if !idCourse.IsZero() {
			if crv, ert := pkg.CRV.GetById(idCourse); ert != nil {
				err = pkg.Trace(ert)
			} else {
				course = crv
			}
		} else if link != "" {
			if crv, ert := pkg.CRV.GetByLink(link); ert != nil {
				err = pkg.Trace(ert)
			} else {
				course = crv
			}
		}

		if err == nil {
			err = pkg.BOT.Send(c, true, course.String(), getTunerBtns(course))
		}
	}

	if showMany && err == nil {
		if skip, ert := strconv.Atoi(paramStr); ert == nil {
			if len(parts) > 2 {
				text = parts[2]
			}

			if courses, ert := pkg.CRV.GetCourses(5, int64(skip), course.GetApproved(), text); ert == nil {
				for _, v := range courses {
					if ert = pkg.BOT.Send(c, true, v.String(), getTunerBtns(v)); ert != nil {
						err = pkg.Trace(ert)
						break
					}
					time.Sleep(time.Millisecond * 200)
				}

				if err == nil {
					if len(courses) > 0 {
						skipL := skip - 5
						skipR := skip + 5

						if skipL < 0 {
							skipL = 0
						}
						inlineMenu := &telebot.ReplyMarkup{}
						inlineMenu.Inline([]telebot.Row{
							{
								telebot.Btn{
									Text:   "<-",
									Unique: fmt.Sprintf("get_approved_courses:%d:%s", skipL, text),
								},
								telebot.Btn{
									Text:   "->",
									Unique: fmt.Sprintf("get_approved_courses:%d:%s", skipR, text),
								},
							},
						}...)

						err = pkg.BOT.Send(c, true, "Iнша сторінка", inlineMenu)

					} else {
						skipL := skip - 5
						if skipL < 0 {
							skipL = 0
						}

						inlineMenu := &telebot.ReplyMarkup{}
						inlineMenu.Inline([]telebot.Row{
							{
								telebot.Btn{
									Text:   "<-",
									Unique: fmt.Sprintf("get_approved_courses:%d:%s", skipL, text),
								},
							},
						}...)

						err = pkg.BOT.Send(c, true, "Курсiв не знайдено\nСпробуйте пошук не пiдтверджених курсiв", inlineMenu)
					}
				}
			} else {
				err = pkg.Trace(ert)
			}
		} else {
			err = pkg.Trace(ert)
		}
	}

	if err != nil {
		return pkg.Trace(err)
	}

	return nil
}
