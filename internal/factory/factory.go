package factory

import (
	"cbot/internal/app"
	"cbot/internal/core"
	"cbot/internal/core/tgcore"
	"cbot/pkg"
)

type FactoryImpl struct{}

func CreateFactory() pkg.Factory {
	return &FactoryImpl{}
}

func (f *FactoryImpl) CreateUserVault() pkg.UserVault {
	return core.CreateUserVault()
}

func (f *FactoryImpl) CreateBot() pkg.Bot {
	return core.CreateBot()
}

func (f *FactoryImpl) CreateTGApp() pkg.TGApp {
	return app.CreateTGApp()
}

func (f *FactoryImpl) CreateCommandVault() pkg.CommandVault {
	return tgcore.CreateCommandVault()
}

func (f *FactoryImpl) CreateUser() pkg.User {
	return core.CreateUser()
}

func (f *FactoryImpl) CreateCourseVault() pkg.CourseVault {
	return core.CreateCourseVault()
}

func (f *FactoryImpl) CreateCourse() pkg.Course {
	return core.CreateCourse()
}

func (f *FactoryImpl) CreateRequestManager() pkg.RequestManager {
	return core.CreateRequestManager()
}

func (f *FactoryImpl) CreateCourseParser() pkg.CourseParser {
	return core.CreateCourseParser()	
}
