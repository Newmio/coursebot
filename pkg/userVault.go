package pkg

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserVault interface {
	CreateOrUpdate(user User) error
	Get(userId int64) (User, error)
	GetAdmins() ([]User, error)
}

type User interface {
	GetObjectId() primitive.ObjectID
	GetId() int64
	GetLogin() string
	GetFirstName() string
	GetMiddleName() string
	GetLastName() string
	GetIsAdmin() bool
	GetGroupNum() int
	GetProffesion() string
	GetProffesionNum() int

	SetObjectId(id primitive.ObjectID)
	SetId(id int64)
	SetLogin(login string)
	SetFirstName(firstName string)
	SetMiddleName(middleName string)
	SetLastName(lastName string)
	SetIsAdmin(isAdmin bool)
	SetGroupNum(groupNum int)
	SetProffesion(proffesion string)
	SetProffesionNum(proffesionNum int)

	String() string
	ToMap() map[string]interface{}
}
