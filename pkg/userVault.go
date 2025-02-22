package pkg

type UserVault interface {
	CreateOrUpdate(user User) error
	Get(userId int64) (User, error)
}

type User interface {
	GetId() int64
	GetLogin() string
	GetFirstName() string 
	GetMiddleName() string
	GetLastName() string

	SetId(id int64)
	SetLogin(login string)
	SetFirstName(firstName string)
	SetMiddleName(middleName string)
	SetLastName(lastName string)

	ToMap() map[string]interface{}
}
