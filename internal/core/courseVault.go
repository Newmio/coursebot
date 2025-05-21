package core

import (
	"bytes"
	"cbot/pkg"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CourseValultImpl struct {
	mClient      *mongo.Collection
	mUserCourses *mongo.Collection
}

func CreateCourseVault() pkg.CourseVault {
	obj := &CourseValultImpl{}
	obj.mClient = pkg.GetMongoCollection(pkg.DBName, pkg.CollectionCourseVault)
	obj.mUserCourses = pkg.GetMongoCollection(pkg.DBName, pkg.CollectionUserCoursesVault)
	return obj
}

func (obj *CourseValultImpl) SetResultCoins(courseId, userId primitive.ObjectID, coins int) error {
	filter := bson.M{"course_id": courseId, "user_id": userId}
	update := bson.M{"$set": bson.M{"result_coins": coins, "coins": bson.A{}}}

	_, err := obj.mUserCourses.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return pkg.Trace(err)
	}

	return nil
}

func (obj *CourseValultImpl) GetCoins(courseId, userId primitive.ObjectID) ([]int, error) {
	filter := bson.M{"course_id": courseId, "user_id": userId}

	var resp bson.M

	if err := obj.mUserCourses.FindOne(context.Background(), filter).Decode(&resp); err != nil {
		return nil, pkg.Trace(err)
	}

	var coins []int

	if valIf, ok := resp["coins"]; ok {
		if valSl, ok := valIf.(primitive.A); ok {
			for _, v := range valSl {
				if val, ok := v.(int32); ok {
					coins = append(coins, int(val))
				}
			}
		}
	}

	return coins, nil
}

func (obj *CourseValultImpl) SetCoins(courseId, userId primitive.ObjectID, coins int) error {
	filter := bson.M{"course_id": courseId, "user_id": userId}
	update := bson.M{"$push": bson.M{"coins": coins}}

	_, err := obj.mUserCourses.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return pkg.Trace(err)
	}

	return nil
}

func (obj *CourseValultImpl) UpdateCheckAdmin(courseId primitive.ObjectID, userId primitive.ObjectID, flag bool) error {
	filter := bson.M{"course_id": courseId, "user_id": userId}
	update := bson.M{"$set": bson.M{"check_admin": true}}

	_, err := obj.mUserCourses.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return pkg.Trace(err)
	}

	return nil
}

func (obj *CourseValultImpl) GetFullFileName(courseId, userId primitive.ObjectID, filename string) (string, error) {
	filter := bson.M{"course_id": courseId, "user_id": userId}

	var resp bson.M

	if err := obj.mUserCourses.FindOne(context.Background(), filter).Decode(&resp); err != nil {
		return "", pkg.Trace(err)
	}

	if valIf, ok := resp["files"]; ok {
		if val, ok := valIf.(primitive.A); ok {
			for _, v := range val {
				str := v.(string)
				if strings.Contains(str, filename) {
					return str, nil
				}
			}
		}
	}

	return "", nil
}

func (obj *CourseValultImpl) DeleteResultFile(courseId, userId primitive.ObjectID, fileName string) error {
	filter := bson.M{"course_id": courseId, "user_id": userId}
	update := bson.M{"$pull": bson.M{"files": fileName}}

	_, err := obj.mUserCourses.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return pkg.Trace(err)
	}

	return nil
}

func (obj *CourseValultImpl) SetResultFile(courseId, userId primitive.ObjectID, fileName string) error {
	filter := bson.M{"course_id": courseId, "user_id": userId}
	update := bson.M{"$push": bson.M{"files": fileName}}

	_, err := obj.mUserCourses.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return pkg.Trace(err)
	}

	return nil
}

func (obj *CourseValultImpl) CheckStartedCourse(courseId, userId primitive.ObjectID) (bool, error) {
	filter := bson.M{"course_id": courseId, "user_id": userId, "start": true, "stop": false}
	count, err := obj.mUserCourses.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, pkg.Trace(err)
	}

	return count > 0, nil
}

func (obj *CourseValultImpl) GetMyCourses(userId primitive.ObjectID) ([]pkg.Course, error) {
	filter := bson.M{"user_id": userId}
	opts := options.Find()

	ctx := context.Background()

	cursor, err := obj.mUserCourses.Find(ctx, filter, opts)
	if err != nil {
		return nil, pkg.Trace(err)
	}

	var courseIds []primitive.ObjectID

	for cursor.Next(ctx) {
		var resp bson.M

		if err := cursor.Decode(&resp); err != nil {
			return nil, pkg.Trace(err)
		}

		if courseIdIf, ok := resp["course_id"]; ok {
			if courseId, ok := courseIdIf.(primitive.ObjectID); ok {
				courseIds = append(courseIds, courseId)
			}
		}
	}

	if len(courseIds) == 0 {
		return nil, nil
	}

	filter = bson.M{"_id": bson.M{"$in": courseIds}}
	opts = options.Find()

	cursor, err = obj.mClient.Find(ctx, filter, opts)
	if err != nil {
		return nil, pkg.Trace(err)
	}

	var courses []pkg.Course

	for cursor.Next(ctx) {
		course := &CourseImpl{}

		var resp bson.M
		if err := cursor.Decode(&resp); err != nil {
			return nil, pkg.Trace(err)
		}

		course.ParseBson(resp)
		courses = append(courses, course)
	}

	return courses, nil
}

func (obj *CourseValultImpl) StartCourse(courseId, userId primitive.ObjectID) error {
	setMap := bson.M{
		"course_id":  courseId,
		"user_id":    userId,
		"start":      true,
		"stop":       false,
		"time_start": time.Now().UTC(),
	}

	opts := options.Update().SetUpsert(true)
	update := bson.M{"$set": setMap}
	filter := bson.M{"course_id": courseId, "user_id": userId}

	_, err := obj.mUserCourses.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return pkg.Trace(err)
	}

	return nil
}

func (obj *CourseValultImpl) StopCourse(courseId, userId primitive.ObjectID) error {
	setMap := bson.M{
		"course_id": courseId,
		"user_id":   userId,
		"start":     false,
		"stop":      true,
		"time_stop": time.Now().UTC(),
	}

	opts := options.Update()
	update := bson.M{"$set": setMap}
	filter := bson.M{"course_id": courseId, "user_id": userId}

	_, err := obj.mUserCourses.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return pkg.Trace(err)
	}

	return nil
}

func (obj *CourseValultImpl) GetCourses(limit, skip int64, appr bool, search string) ([]pkg.Course, error) {
	filter := make(bson.M)
	opts := options.Find().SetLimit(limit).SetSkip(skip)

	if appr {
		filter["approved"] = true
	}

	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"description": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	courses := make([]pkg.Course, 0)
	ctx := context.Background()

	cursor, err := obj.mClient.Find(ctx, filter, opts)
	if err != nil {
		return nil, pkg.Trace(err)
	}

	for cursor.Next(ctx) {
		course := &CourseImpl{}

		var resp bson.M
		if err := cursor.Decode(&resp); err != nil {
			return nil, pkg.Trace(err)
		}

		course.ParseBson(resp)
		courses = append(courses, course)
	}

	return courses, nil
}

func (obj *CourseValultImpl) CreateOrUpdate(course pkg.Course) error {
	var filter bson.M

	opts := options.Update().SetUpsert(true)
	update := bson.M{"$set": course.ToMap()}

	if !course.GetId().IsZero() {
		filter = bson.M{"_id": course.GetId()}
	} else if course.GetLink() != "" {
		filter = bson.M{"link": course.GetLink()}
	} else {
		return pkg.Trace(fmt.Errorf("empty course id and link"))
	}

	_, err := obj.mClient.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return pkg.Trace(err)
	}

	return nil
}

func (obj *CourseValultImpl) GetByLink(link string) (pkg.Course, error) {
	filter := bson.M{"link": link}
	course := &CourseImpl{}
	var resp bson.M

	if err := obj.mClient.FindOne(context.Background(), filter).Decode(&resp); err != nil {
		return nil, pkg.Trace(err)
	}

	course.ParseBson(resp)

	return course, nil
}

func (obj *CourseValultImpl) GetById(id primitive.ObjectID) (pkg.Course, error) {
	filter := bson.M{"_id": id}
	course := &CourseImpl{}
	var resp bson.M

	if err := obj.mClient.FindOne(context.Background(), filter).Decode(&resp); err != nil {
		return nil, pkg.Trace(err)
	}

	course.ParseBson(resp)

	return course, nil
}

func (obj *CourseValultImpl) Exists(link string) (bool, error) {
	filter := bson.M{"link": link}
	count, err := obj.mClient.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, pkg.Trace(err)
	}

	return count > 0, nil
}

type CourseImpl struct {
	mId          primitive.ObjectID
	mName        string
	mDescription string
	mCost        string
	mDuration    string
	mApproved    bool
	mLink        string
}

func CreateCourse() pkg.Course {
	return &CourseImpl{}
}

func (obj *CourseImpl) String() string {
	if obj.GetName() == "" {
		obj.SetName("-")
	}

	if obj.GetDescription() == "" {
		obj.SetDescription("-")
	}

	if obj.GetDuration() == "" {
		obj.SetDuration("-")
	}

	if obj.GetLink() == "" {
		obj.SetLink("-")
	}

	return fmt.Sprintf("Назва: %s\nОпис: %s\nЦіна: %s\nТривалість: %s\nПосилання: %s",
		obj.GetName(), obj.GetDescription(), obj.GetCost(), obj.GetDuration(), obj.GetLink())
}

func (obj *CourseImpl) GetId() primitive.ObjectID {
	return obj.mId
}

func (obj *CourseImpl) GetName() string {
	return obj.mName
}

func (obj *CourseImpl) GetDescription() string {
	return obj.mDescription
}

func (obj *CourseImpl) GetCost() string {
	return obj.mCost
}

func (obj *CourseImpl) GetDuration() string {
	return obj.mDuration
}

func (obj *CourseImpl) GetApproved() bool {
	return obj.mApproved
}

func (obj *CourseImpl) GetLink() string {
	return obj.mLink
}

func (obj *CourseImpl) SetId(id primitive.ObjectID) {
	obj.mId = id
}

func (obj *CourseImpl) SetName(name string) {
	obj.mName = name
}

func (obj *CourseImpl) SetDescription(description string) {
	obj.mDescription = description
}

func (obj *CourseImpl) SetCost(cost string) {
	obj.mCost = cost
}

func (obj *CourseImpl) SetDuration(duration string) {
	obj.mDuration = duration
}

func (obj *CourseImpl) SetApproved(approved bool) {
	obj.mApproved = approved
}

func (obj *CourseImpl) SetLink(link string) {
	obj.mLink = link
}

func (obj *CourseImpl) ParseBson(bsonM bson.M) {
	if idIf, ok := bsonM["_id"]; ok {
		if id, ok := idIf.(primitive.ObjectID); ok {
			obj.SetId(id)
		}
	}

	if nameIf, ok := bsonM["name"]; ok {
		if name, ok := nameIf.(string); ok {
			obj.SetName(name)
		}
	}

	if descriptionIf, ok := bsonM["description"]; ok {
		if description, ok := descriptionIf.(string); ok {
			obj.SetDescription(description)
		}
	}

	if costIf, ok := bsonM["cost"]; ok {
		if cost, ok := costIf.(string); ok {
			obj.SetCost(cost)
		}
	}

	if durationIf, ok := bsonM["duration"]; ok {
		if duration, ok := durationIf.(string); ok {
			obj.SetDuration(duration)
		}
	}

	if approvedIf, ok := bsonM["approved"]; ok {
		if approved, ok := approvedIf.(bool); ok {
			obj.SetApproved(approved)
		}
	}

	if linkIf, ok := bsonM["link"]; ok {
		if link, ok := linkIf.(string); ok {
			obj.SetLink(link)
		}
	}
}

func (obj *CourseImpl) ToMap() map[string]interface{} {
	resp := make(map[string]interface{})

	if obj.mId != primitive.NilObjectID {
		resp["_id"] = obj.mId
	}

	if obj.mName != "" {
		resp["name"] = obj.mName
	}

	if obj.mDescription != "" {
		resp["description"] = obj.mDescription
	}

	if obj.mCost != "" {
		resp["cost"] = obj.mCost
	}

	if obj.mDuration != "" {
		resp["duration"] = obj.mDuration
	}

	if obj.mLink != "" {
		resp["link"] = obj.mLink
	}

	resp["approved"] = obj.mApproved

	return resp
}

type CourseParserImpl struct {
}

func CreateCourseParser() pkg.CourseParser {
	return &CourseParserImpl{}
}

func (obj *CourseParserImpl) StartParseSite(searchValue string, siteName string) ([]map[string]string, error) {
	var params map[string]string
	var siteLink string

	siteParamsIf, ok := pkg.CoursesParameters[siteName]
	if !ok {
		return nil, pkg.Trace(fmt.Errorf("site %s not found", siteName))
	}

	siteParamsMap, ok := siteParamsIf.(map[string]interface{})
	if !ok {
		return nil, pkg.Trace(fmt.Errorf("invalid type for site %s", siteName))
	}

	siteLinkIf, ok := siteParamsMap["site_link"]
	if !ok {
		return nil, pkg.Trace(fmt.Errorf("site link not found in site %s", siteName))
	}

	siteLink, ok = siteLinkIf.(string)
	if !ok {
		return nil, pkg.Trace(fmt.Errorf("invalid type for site link in site %s", siteName))
	}

	paramsFealdsIf, ok := siteParamsMap["fealds"]
	if !ok {
		return nil, pkg.Trace(fmt.Errorf("fealds not found in site %s", siteName))
	}

	params, ok = paramsFealdsIf.(map[string]string)
	if !ok {
		return nil, pkg.Trace(fmt.Errorf("invalid type for fealds in site %s", siteName))
	}

	var pagination bool
	if strings.Contains(siteLink, "<page>") {
		pagination = true
	}

	ret := make([]map[string]string, 0)

	page := 0
	restartReq := true
	for {
		siteLink = strings.ReplaceAll(siteLink, "<search_value>", searchValue)
		siteLink = strings.ReplaceAll(siteLink, "<page>", strconv.Itoa(page))

		req, err := http.NewRequest("GET", siteLink, nil)
		if err != nil {
			return nil, pkg.Trace(err)
		}

		parsedURL, err := url.ParseRequestURI(siteLink)
		if err != nil {
			return nil, pkg.Trace(err)
		}

		req.Header.Set("Host", parsedURL.Host)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:107.0) Gecko/20100101 Firefox/107.0")

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			return nil, pkg.Trace(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, pkg.Trace(err)
		}

		if resp.StatusCode != 200 {
			if !restartReq {
				return nil, fmt.Errorf("error request to %s, code %d", siteLink, resp.StatusCode)
			} else {
				time.Sleep(time.Second)
				restartReq = false
				continue
			}
		}
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			return nil, pkg.Trace(err)
		}

		retMap := obj.ParseParams(params, doc)

		if retMap != nil {
			ret = append(ret, retMap)
		} else {
			break
		}

		if !pagination || page > 15 {
			break
		}

		page++
	}

	return ret, nil
}

func (obj *CourseParserImpl) ParseParams(params map[string]string, doc *goquery.Document) map[string]string {
	resp := make(map[string]string)

	//node := doc.Find(params["main"])

	for key, value := range params {
		var selector, attr string

		if key == "main" {
			continue
		}

		parts := strings.Split(value, "<>")
		if len(parts) == 1 {
			selector = parts[0]
		} else {
			selector = parts[0]
			attr = parts[1]
		}

		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			if attr != "" {
				if str, ok := s.Attr(attr); ok {
					resp[key] = str
				} else {
					resp[key] = fmt.Sprintf("%s not found in %s", attr, selector)
				}
			} else {
				resp[key] = strings.ReplaceAll(s.Text(), "\n", "")
			}
		})
	}

	return resp
}
