package tgcore

import (
	"cbot/pkg"
	"sync"

	"gopkg.in/telebot.v4"
)

type CommandVaultImpl struct {
	mMu       sync.Mutex
	mCommands map[int64]string

	mDeleteMessages map[int64][]*telebot.Message
}

func CreateCommandVault() pkg.CommandVault {
	obj := &CommandVaultImpl{}
	obj.mCommands = make(map[int64]string)
	obj.mDeleteMessages = make(map[int64][]*telebot.Message)
	return obj
}

func (obj *CommandVaultImpl) SetCommand(userId int64, command string) {
	obj.mMu.Lock()
	defer obj.mMu.Unlock()

	obj.mCommands[userId] = command
}

func (obj *CommandVaultImpl) GetCommand(userId int64) string {
	return obj.mCommands[userId]
}

func (obj *CommandVaultImpl) RemoveCommand(userId int64) {
	obj.mMu.Lock()
	defer obj.mMu.Unlock()

	delete(obj.mCommands, userId)
}

func (obj *CommandVaultImpl) AppendDeleteMessage(msg *telebot.Message) {
	obj.mMu.Lock()
	defer obj.mMu.Unlock()

	obj.mDeleteMessages[msg.Chat.ID] = append(obj.mDeleteMessages[msg.Chat.ID], msg)
}

func (obj *CommandVaultImpl) ClearAllDeleteMessages(userId int64) {
	if messages, ok := obj.mDeleteMessages[userId]; ok && len(messages) > 0 {
		editables := make([]telebot.Editable, 0, len(messages))

		for _, msg := range messages {
			editables = append(editables, msg)
		}

		pkg.BOT.GetBot().DeleteMany(editables)

		obj.mMu.Lock()
		defer obj.mMu.Unlock()

		obj.mDeleteMessages = make(map[int64][]*telebot.Message)
	}
}
