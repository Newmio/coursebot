package core

import (
	"cbot/pkg"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserVaultImpl struct {
	mClient *mongo.Collection
}

func CreateUserVault() pkg.UserVault {
	obj := &UserVaultImpl{}
	obj.mClient = pkg.GetMongoCollection(pkg.DBName, pkg.CollectionUserVault)
	return obj
}

func (obj *UserVaultImpl) CreateOrUpdate(user pkg.User) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"id": user.GetId()}
	update := bson.M{"$set": user.ToMap()}

	_, err := obj.mClient.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return pkg.Trace(err)
	}

	return nil
}

func (obj *UserVaultImpl) Get(userId int64) (pkg.User, error) {
	filter := bson.M{"id": userId}
	user := &UserImpl{}

	var resp bson.M

	if err := obj.mClient.FindOne(context.Background(), filter).Decode(&resp); err != nil {
		return nil, pkg.Trace(err)
	}

	if idIf, ok := resp["id"]; ok {
		if id, ok := idIf.(int64); ok {
			user.SetId(id)
		}
	}

	if loginIf, ok := resp["login"]; ok {
		if login, ok := loginIf.(string); ok {
			user.SetLogin(login)
		}
	}

	if firstNameIf, ok := resp["first_name"]; ok {
		if firstName, ok := firstNameIf.(string); ok {
			user.SetFirstName(firstName)
		}
	}

	if middleNameIf, ok := resp["middle_name"]; ok {
		if middleName, ok := middleNameIf.(string); ok {
			user.SetMiddleName(middleName)
		}
	}

	if lastNameIf, ok := resp["last_name"]; ok {
		if lastName, ok := lastNameIf.(string); ok {
			user.SetLastName(lastName)
		}
	}

	return user, nil
}

type UserImpl struct {
	mId         int64
	mLogin      string
	mFirstName  string
	mMiddleName string
	mLastName   string
}

func CreateUser() pkg.User {
	return &UserImpl{}
}

func (obj *UserImpl) GetId() int64 {
	return obj.mId
}

func (obj *UserImpl) GetLogin() string {
	return obj.mLogin
}

func (obj *UserImpl) GetFirstName() string {
	return obj.mFirstName
}

func (obj *UserImpl) GetMiddleName() string {
	return obj.mMiddleName
}

func (obj *UserImpl) GetLastName() string {
	return obj.mLastName
}

func (obj *UserImpl) SetId(id int64) {
	obj.mId = id
}

func (obj *UserImpl) SetLogin(login string) {
	obj.mLogin = login
}

func (obj *UserImpl) SetFirstName(firstName string) {
	obj.mFirstName = firstName
}

func (obj *UserImpl) SetMiddleName(middleName string) {
	obj.mMiddleName = middleName
}

func (obj *UserImpl) SetLastName(lastName string) {
	obj.mLastName = lastName
}

func (obj *UserImpl) ToMap() map[string]interface{} {
	resp := make(map[string]interface{})

	if obj.mId != 0 {
		resp["id"] = obj.mId
	}

	if obj.mLogin != "" {
		resp["login"] = obj.mLogin
	}

	if obj.mFirstName != "" {
		resp["first_name"] = obj.mFirstName
	}

	if obj.mMiddleName != "" {
		resp["middle_name"] = obj.mMiddleName
	}

	if obj.mLastName != "" {
		resp["last_name"] = obj.mLastName
	}

	return resp
}