package framework

import (
	"context"

	"github.com/disgoorg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func MongoConnect(uri string) {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Error("Failed to connect to MongoDB.")
	} else {
		log.Info("Connected to MongoDB.")
	}
}

func FindDocument(collection string, filter bson.D) bson.M {
	coll := client.Database("database").Collection(collection)
	var result bson.M
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		panic(err)
	}
	return result
}

func FindDocuments(collection string, filter bson.D) []bson.M {
	coll := client.Database("database").Collection(collection)
	var results []bson.M
	cursor, err := coll.Find(context.TODO(), filter)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		panic(err)
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return results
}

func CountDocuments(collection string, filter bson.D) int64 {
	coll := client.Database("database").Collection(collection)
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	return count
}

func InsertDocument(collection string, document bson.D) *mongo.InsertOneResult {
	coll := client.Database("database").Collection(collection)
	result, err := coll.InsertOne(context.TODO(), document)
	if err != nil {
		panic(err)
	}
	return result
}

func InsertDocuments(collection string, documents []interface{}) *mongo.InsertManyResult {
	coll := client.Database("database").Collection(collection)
	result, err := coll.InsertMany(context.TODO(), documents)
	if err != nil {
		panic(err)
	}
	return result
}

func UpdateDocument(collection string, filter bson.D, update bson.D) *mongo.UpdateResult {
	coll := client.Database("database").Collection(collection)
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return result
}

func UpdateDocuments(collection string, filter bson.D, update bson.D) *mongo.UpdateResult {
	coll := client.Database("database").Collection(collection)
	result, err := coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return result
}

func DeleteDocument(collection string, filter bson.D) *mongo.DeleteResult {
	coll := client.Database("database").Collection(collection)
	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	return result
}
func DeleteDocuments(collection string, filter bson.D) *mongo.DeleteResult {
	coll := client.Database("database").Collection(collection)
	result, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	return result
}

type Infraction struct {
	ID          int
	UserID      string
	GuildID     string
	Reason      string
	Type        string
	ModeratorID string
}

type Guild struct {
	ID       string
	MuteRole string
}

func (infraction Infraction) Data() bson.D {
	return bson.D{
		{Key: "ID", Value: infraction.ID},
		{Key: "UserID", Value: infraction.UserID},
		{Key: "ModeratorID", Value: infraction.ModeratorID},
		{Key: "GuildID", Value: infraction.GuildID},
		{Key: "Reason", Value: infraction.Reason},
		{Key: "Type", Value: infraction.Type},
	}
}

func (guild Guild) Data() bson.D {
	return bson.D{
		{Key: "ID", Value: guild.ID},
		{Key: "MuteRole", Value: guild.MuteRole},
	}
}
