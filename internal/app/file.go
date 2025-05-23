package app

import (
	"cbot/pkg"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/telebot.v4"
)

func (obj *TGAppImpl) getCourseFiles(c telebot.Context, courseIdStr, userIdStr string) error {
	courseId, err := primitive.ObjectIDFromHex(courseIdStr)
	if err != nil {
		return pkg.BOT.Send(c, false, pkg.Trace(err).Error())
	}

	userId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		return pkg.BOT.Send(c, false, pkg.Trace(err).Error())
	}

	courseInfo, err := pkg.CRV.GetUserCourse(courseId, userId)
	if err != nil {
		return pkg.BOT.Send(c, false, pkg.Trace(err).Error())
	}

	if len(courseInfo) > 0 {
		if filesIf, ok := courseInfo["files"]; ok {
			if files, ok := filesIf.(primitive.A); ok && len(files) > 0 {
				for _, v := range files {
					if str, ok := v.(string); ok {

						file, err := pkg.FLV.Get(str)
						if err != nil {
							return pkg.BOT.Send(c, false, pkg.Trace(err).Error())
						}
						defer file.Close()

						inlineMenu := &telebot.ReplyMarkup{}
						inlineMenu.Inline([]telebot.Row{
							{
								telebot.Btn{
									Text:   "Так",
									Unique: fmt.Sprintf("btn_d_c_f:%s:%s", courseIdStr, strings.Split(str, ".")[0]),
								},
								telebot.Btn{
									Text:   "Ні",
									Unique: "btn_clear_msg",
								},
							},
						}...)

						if strings.Contains(str, "jpg") {
							photo := &telebot.Photo{File: telebot.FromReader(file), Caption: "Видалити файл?"}
							if err := pkg.BOT.Send(c, true, photo, inlineMenu); err != nil {
								return pkg.BOT.Send(c, false, pkg.Trace(err).Error())
							}
						} else {
							document := &telebot.Document{File: telebot.FromReader(file), Caption: "Видалити файл?"}
							if err := pkg.BOT.Send(c, true, document, inlineMenu); err != nil {
								return pkg.BOT.Send(c, false, pkg.Trace(err).Error())
							}
						}

						file.Close()
					}
				}
			} else {
				return pkg.BOT.Send(c, true, "Файли відсутні")
			}
		}
	}

	return nil
}

func (obj *TGAppImpl) deleteFilesByCourse(c telebot.Context, courseIdStr, fileName string) error {
	user := c.Get("user").(pkg.User)

	courseId, err := primitive.ObjectIDFromHex(courseIdStr)
	if err != nil {
		return pkg.BOT.Send(c, false, pkg.Trace(err).Error())
	}

	fullFileName, err := pkg.CRV.GetFullFileName(courseId, user.GetObjectId(), fileName)
	if err != nil {
		return pkg.BOT.Send(c, false, pkg.Trace(err).Error())
	}

	if err := pkg.CRV.DeleteResultFile(courseId, user.GetObjectId(), fullFileName); err != nil {
		return pkg.BOT.Send(c, false, pkg.Trace(err).Error())
	}

	if err := pkg.FLV.Delete(fullFileName); err != nil {
		return pkg.BOT.Send(c, false, pkg.Trace(err).Error())
	}

	return pkg.BOT.Send(c, true, "Файл видалено")
}

func (obj *TGAppImpl) courseResultFileHandler(c telebot.Context, courseId primitive.ObjectID) error {
	user := c.Get("user").(pkg.User)
	msg := c.Message()
	var file *telebot.File

	fileName := fmt.Sprint(user.GetId()) + "/" + strconv.Itoa(msg.ID)

	if msg.Document != nil {
		file = &msg.Document.File
		fileName += filepath.Ext(msg.Document.FileName)

	} else if msg.Photo != nil {
		file = &msg.Photo.File
		fileName += ".jpg"

	} else {
		return pkg.BOT.Send(c, false, "Невідомий тип повідомлення")
	}

	if err := pkg.FLV.Create(nil, fileName); err != nil {
		return pkg.BOT.Send(c, false, pkg.Trace(err).Error())
	}

	if err := pkg.BOT.GetBot().Download(file, pkg.FLV.GetDir(fileName)); err != nil {
		return pkg.BOT.Send(c, false, pkg.Trace(err).Error())
	}

	if err := pkg.CRV.SetResultFile(courseId, user.GetObjectId(), fileName); err != nil {
		return pkg.BOT.Send(c, false, pkg.Trace(err).Error())
	}

	idStr := courseId.Hex()

	inlineMenu := &telebot.ReplyMarkup{}

	inlineMenu.Inline([]telebot.Row{
		{
			telebot.Btn{
				Text:   "Так",
				Unique: fmt.Sprintf("btn_d_c_f:%s:%s", idStr, strings.Split(fileName, ".")[0]),
			},
			telebot.Btn{
				Text:   "Ні",
				Unique: "btn_clear_msg",
			},
		},
	}...)

	if msg.Document != nil {
		msg.Document.Caption = "Документ успішно завантажено\nВидалити?"
		err := pkg.BOT.Send(c, false, msg.Document, inlineMenu)
		if err != nil {
			return pkg.BOT.Send(c, true, pkg.Trace(err).Error())
		}

	} else if msg.Photo != nil {
		msg.Photo.Caption = "Фото успішно завантажено\nВидалити?"
		err := pkg.BOT.Send(c, true, msg.Photo, inlineMenu)
		if err != nil {
			return pkg.BOT.Send(c, false, pkg.Trace(err).Error())
		}

	} else {
		return pkg.BOT.Send(c, false, "Невідомий тип повідомлення")
	}

	inlineMenu.Inline([]telebot.Row{
		{
			telebot.Btn{
				Text:   "Вiдправити результати",
				Unique: fmt.Sprintf("btn_complete_send_result_course:%s", idStr),
			},
		},
	}...)

	pkg.CMDV.SetCommand(c.Sender().ID, fmt.Sprintf("set_course_result:%s", idStr))

	return pkg.BOT.Send(c, false, "Вiдправте результат роботи\nАбо надішліть ще файлiв", inlineMenu)
}
