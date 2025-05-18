package pkg

type Factory interface {
	CreateUserVault() UserVault
	CreateBot() Bot
	CreateTGApp() TGApp
	CreateCommandVault() CommandVault
	CreateUser() User
	CreateCourseVault() CourseVault
	CreateCourse() Course
	CreateRequestManager() RequestManager
	CreateCourseParser() CourseParser
	CreateFileVault() FileVault
}
