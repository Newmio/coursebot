package app

import (
	"cbot/pkg"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/telebot.v4"
)

func (obj *TGAppImpl) courseResultFileHandler(c telebot.Context, courseId primitive.ObjectID) error {
	user := c.Get("user").(pkg.User)
	msg := c.Message()
	var file *telebot.File
	fileName := fmt.Sprint(user.GetId()) + "/"

	if msg.Document != nil {
		file = &msg.Document.File
		fileName += msg.Document.FileName
	} else if msg.Photo != nil {
		file = &msg.Photo.File
		fileName += strconv.Itoa(msg.ID) + ".jpg"
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

	return pkg.BOT.Send(c, true, "Файл успішно завантажено")
}
