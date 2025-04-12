package main

import (
	"cbot/internal/factory"
	"cbot/pkg"
	"os"
)

func main() {
	pkg.MongoHost = os.Getenv("MONGO_HOST")
	pkg.BotToken = os.Getenv("TG_TOKEN")
	//owner := os.Getenv("TG_OWNER")

	pkg.F = factory.CreateFactory()
	pkg.BOT = pkg.F.CreateBot()
	pkg.USRV = pkg.F.CreateUserVault()
	pkg.CMDV = pkg.F.CreateCommandVault()
	pkg.CRV = pkg.F.CreateCourseVault()

	pkg.F.CreateTGApp().Run()
}
