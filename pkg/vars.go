package pkg

var (
	F Factory

	// Mongo
	USRV UserVault

	// Telegram
	BOT  Bot
	CMDV CommandVault

	BotToken  string
	MongoHost string
)

const (
	DBName                = "cbot"
	CollectionUserVault   = "users"
	CollectionCourseVault = "courses"
)
