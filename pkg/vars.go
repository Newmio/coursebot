package pkg

var (
	F Factory

	// Mongo
	USRV UserVault
	CRV  CourseVault

	// Telegram
	BOT  Bot
	CMDV CommandVault

	BotToken  string
	MongoHost string

	CoursesParameters = map[string]interface{}{
		"prometheus": map[string]interface{}{
			"site_link": "https://prometheus.org.ua/courses-catalog?q=<search_value>",
			"fealds": map[string]string{
				"main": "div.course-list.flex.flex-col.items-center.md\\:grid.md\\:grid-cols-3.md\\:auto-rows-max.lg\\:grid-cols-2.3colCatalog\\:grid-cols-3.gap-12.md\\:gap-8.w-full",
				"link": "a<>href",
			},
		},
	}
)

const (
	DBName                     = "cbot"
	CollectionUserVault        = "users"
	CollectionCourseVault      = "courses"
	CollectionUserCoursesVault = "user_courses"
)
