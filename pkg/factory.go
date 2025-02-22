package pkg

type Factory interface {
	CreateUserVault() UserVault
	CreateBot() Bot
	CreateTGApp() TGApp
	CreateCommandVault() CommandVault
	CreateUser() User
}
