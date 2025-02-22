package pkg

type CommandVault interface {
	SetCommand(userId int64, command string)
	GetCommand(userId int64) string
	RemoveCommand(userId int64)
}
