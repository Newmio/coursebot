package tgcore

import (
	"cbot/pkg"
	"sync"
)

type CommandVaultImpl struct {
	mMu sync.Mutex
	mCommands map[int64]string
}

func CreateCommandVault() pkg.CommandVault{
	obj := &CommandVaultImpl{}
	obj.mCommands = make(map[int64]string)
	return obj
}

func (obj *CommandVaultImpl) SetCommand(userId int64, command string)  {
	obj.mMu.Lock()
	defer obj.mMu.Unlock()

	obj.mCommands[userId] = command
}

func (obj *CommandVaultImpl) GetCommand(userId int64) string  {
	return obj.mCommands[userId]
}

func (obj *CommandVaultImpl) RemoveCommand(userId int64)  {
	obj.mMu.Lock()
	defer obj.mMu.Unlock()

	delete(obj.mCommands, userId)
}