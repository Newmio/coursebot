package pkg

import "go.mongodb.org/mongo-driver/bson/primitive"

type CourseVault interface {
	CreateOrUpdate(course Course) error
}

type Course interface {
	GetId() primitive.ObjectID
	GetName() string
	GetDescription() string
	GetCost() int
	GetDuration() string
	GetApproved() bool
	GetLink() string

	SetId(id primitive.ObjectID)
	SetName(name string)
	SetDescription(description string)
	SetCost(cost int)
	SetDuration(duration string)
	SetApproved(approved bool)
	SetLink(link string)

	ToMap() map[string]interface{}
}


type CourseParser interface{

}