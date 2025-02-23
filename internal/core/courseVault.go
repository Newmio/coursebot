package core

import (
	"cbot/pkg"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CourseValultImpl struct {
	mClient *mongo.Collection
}

func CreateCourseVault() pkg.CourseVault {
	obj := &CourseValultImpl{}
	obj.mClient = pkg.GetMongoCollection(pkg.DBName, pkg.CollectionCourseVault)
	return obj
}

func (obj *CourseValultImpl) CreateOrUpdate(course pkg.Course) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"link": course.GetLink()}
	update := bson.M{"$set": course.ToMap()}

	_, err := obj.mClient.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return pkg.Trace(err)
	}

	return nil
}

type Course struct {
	mId          primitive.ObjectID
	mName        string
	mDescription string
	mCost        int
	mDuration    string
	mApproved    bool
	mLink        string
}

func CreateCourse() pkg.Course {
	return &Course{}
}

func (obj *Course) GetId() primitive.ObjectID {
	return obj.mId
}

func (obj *Course) GetName() string {
	return obj.mName
}

func (obj *Course) GetDescription() string {
	return obj.mDescription
}

func (obj *Course) GetCost() int {
	return obj.mCost
}

func (obj *Course) GetDuration() string {
	return obj.mDuration
}

func (obj *Course) GetApproved() bool {
	return obj.mApproved
}

func (obj *Course) GetLink() string {
	return obj.mLink
}

func (obj *Course) SetId(id primitive.ObjectID) {
	obj.mId = id
}

func (obj *Course) SetName(name string) {
	obj.mName = name
}

func (obj *Course) SetDescription(description string) {
	obj.mDescription = description
}

func (obj *Course) SetCost(cost int) {
	obj.mCost = cost
}

func (obj *Course) SetDuration(duration string) {
	obj.mDuration = duration
}

func (obj *Course) SetApproved(approved bool) {
	obj.mApproved = approved
}

func (obj *Course) SetLink(link string) {
	obj.mLink = link
}

func (obj *Course) ToMap() map[string]interface{} {
	resp := make(map[string]interface{})

	if obj.mId != primitive.NilObjectID {
		resp["id"] = obj.mId
	}

	if obj.mName != "" {
		resp["name"] = obj.mName
	}

	if obj.mDescription != "" {
		resp["description"] = obj.mDescription
	}

	if obj.mCost != 0 {
		resp["cost"] = obj.mCost
	}

	if obj.mDuration != "" {
		resp["duration"] = obj.mDuration
	}

	if obj.mLink != "" {
		resp["link"] = obj.mLink
	}

	return resp
}

type CourseParserImpl struct {
	 
}

func CreateCourseParser() pkg.CourseParser {
	return &CourseParserImpl{}
}