package pkg

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CourseVault interface {
	CreateOrUpdate(course Course) error
	Exists(link string) (bool, error)
	GetByLink(link string) (Course, error)
	GetById(id primitive.ObjectID) (Course, error)
	GetCourses(limit, skip int64, appr bool, search string) ([]Course, error)
	StartCourse(courseId, userId primitive.ObjectID) error
	StopCourse(courseId, userId primitive.ObjectID) error
	GetMyCourses(userId primitive.ObjectID) ([]Course, error)
	CheckStartedCourse(courseId, userId primitive.ObjectID) (bool, error)
	SetResultFile(courseId, userId primitive.ObjectID, fileName string)error
	DeleteResultFile(courseId, userId primitive.ObjectID, fileName string) error
}

type Course interface {
	GetId() primitive.ObjectID
	GetName() string
	GetDescription() string
	GetCost() string
	GetDuration() string
	GetApproved() bool
	GetLink() string

	SetId(id primitive.ObjectID)
	SetName(name string)
	SetDescription(description string)
	SetCost(cost string)
	SetDuration(duration string)
	SetApproved(approved bool)
	SetLink(link string)

	ToMap() map[string]interface{}
	ParseBson(bsonM bson.M)
	String() string
}

type CourseParser interface {
	StartParseSite(searchValue string, siteName string) ([]map[string]string, error)
}
