package pkg

import (
	"context"
	"fmt"
	"runtime"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoCollection(dbName, collection string) *mongo.Collection {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(MongoHost))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		panic(err)
	}

	return client.Database(dbName).Collection(collection)
}

func Trace(err error, any ...interface{}) error {
	if err == nil {
		return nil
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return err
	}

	var str string

	for _, value := range any {

		switch v := value.(type) {
		default:
			str += fmt.Sprint(v)
		}
	}

	return fmt.Errorf("%s%s%s%s(*_*) %s:%d (*_*)", err.Error(), "\n", str, "\n", file, line)
}
