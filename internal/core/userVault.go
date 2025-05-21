package core

import (
	"cbot/pkg"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (obj *UserVaultImpl) GetAdmins() ([]pkg.User, error) {
	filter := bson.M{"is_admin": true}
	cursor, err := obj.mClient.Find(context.Background(), filter)
	if err != nil {
		return nil, pkg.Trace(err)
	}

	var resp bson.M
	users := make([]pkg.User, 0)

	for cursor.Next(context.Background()) {
		user := &UserImpl{}

		if err := cursor.Decode(&resp); err != nil {
			return nil, pkg.Trace(err)
		}

		if objIdIf, ok := resp["_id"]; ok {
			if objId, ok := objIdIf.(primitive.ObjectID); ok {
				user.SetObjectId(objId)
			}
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

		if isAdminIf, ok := resp["is_admin"]; ok {
			if isAdmin, ok := isAdminIf.(bool); ok {
				user.SetIsAdmin(isAdmin)
			}
		}

		if groupNumIf, ok := resp["group_num"]; ok {
			if groupNum, ok := groupNumIf.(int); ok {
				user.GroupNum = groupNum
			}
		}

		if proffesionIf, ok := resp["proffesion"]; ok {
			if proffesion, ok := proffesionIf.(string); ok {
				user.Proffesion = proffesion
			}
		}

		if proffesionNumIf, ok := resp["proffesion_num"]; ok {
			if proffesionNum, ok := proffesionNumIf.(int); ok {
				user.ProffesionNum = proffesionNum
			}
		}

		users = append(users, user)
	}

	return users, nil
}

func (obj *UserVaultImpl) Get(userId int64) (pkg.User, error) {
	filter := bson.M{"id": userId}
	user := &UserImpl{}

	var resp bson.M

	if err := obj.mClient.FindOne(context.Background(), filter).Decode(&resp); err != nil {
		return nil, pkg.Trace(err)
	}

	if objIdIf, ok := resp["_id"]; ok {
		if objId, ok := objIdIf.(primitive.ObjectID); ok {
			user.SetObjectId(objId)
		}
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

	if isAdminIf, ok := resp["is_admin"]; ok {
		if isAdmin, ok := isAdminIf.(bool); ok {
			user.SetIsAdmin(isAdmin)
		}
	}

	if groupNumIf, ok := resp["group_num"]; ok {
		if groupNum, ok := groupNumIf.(int); ok {
			user.GroupNum = groupNum
		}
	}

	if proffesionIf, ok := resp["proffesion"]; ok {
		if proffesion, ok := proffesionIf.(string); ok {
			user.Proffesion = proffesion
		}
	}

	if proffesionNumIf, ok := resp["proffesion_num"]; ok {
		if proffesionNum, ok := proffesionNumIf.(int); ok {
			user.ProffesionNum = proffesionNum
		}
	}

	return user, nil
}

type UserImpl struct {
	mObjectId     primitive.ObjectID
	mId           int64
	mLogin        string
	mFirstName    string
	mMiddleName   string
	mLastName     string
	mIsAdmin      bool
	GroupNum      int
	Proffesion    string
	ProffesionNum int
}

func CreateUser() pkg.User {
	return &UserImpl{}
}

func (obj *UserImpl) String() string {
	text := fmt.Sprintf("Iм'я: %s\n", obj.GetFirstName())
	text += fmt.Sprintf("Прізвище: %s\n", obj.GetLastName())
	text += fmt.Sprintf("Логін: %s\n", obj.GetMiddleName())

	if obj.GetIsAdmin() {
		text += "Роль: Адміністратор\n"
	} else {
		text += "Роль: Студент\n"
	}

	text += fmt.Sprintf("Курс: %d\n", obj.GroupNum)
	text += fmt.Sprintf("Професія: %s\n", obj.Proffesion)
	text += fmt.Sprintf("Номер професії: %d\n", obj.ProffesionNum)

	return text
}

func (obj *UserImpl) GetObjectId() primitive.ObjectID {
	return obj.mObjectId
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

func (obj *UserImpl) GetIsAdmin() bool {
	return obj.mIsAdmin
}

func (obj *UserImpl) GetGroupNum() int {
	return obj.GroupNum
}

func (obj *UserImpl) GetProffesion() string {
	return obj.Proffesion
}

func (obj *UserImpl) GetProffesionNum() int {
	return obj.ProffesionNum
}

func (obj *UserImpl) SetObjectId(id primitive.ObjectID) {
	obj.mObjectId = id
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

func (obj *UserImpl) SetIsAdmin(isAdmin bool) {
	obj.mIsAdmin = isAdmin
}

func (obj *UserImpl) SetGroupNum(groupNum int) {
	obj.GroupNum = groupNum
}

func (obj *UserImpl) SetProffesion(proffesion string) {
	obj.Proffesion = proffesion
}

func (obj *UserImpl) SetProffesionNum(proffesionNum int) {
	obj.ProffesionNum = proffesionNum
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

	if obj.GroupNum != 0 {
		resp["group_num"] = obj.GroupNum
	}

	if obj.Proffesion != "" {
		resp["proffesion"] = obj.Proffesion
	}

	if obj.ProffesionNum != 0 {
		resp["proffesion_num"] = obj.ProffesionNum
	}

	resp["is_admin"] = obj.mIsAdmin

	return resp
}
