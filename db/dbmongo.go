package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//var collection *mongo.Collection
var ctx = context.TODO()

type MongoDb struct {
	mdbInstance *mongo.Database
}

type collectionInfo struct {
	entityName    string
	collectionHdl *mongo.Collection
}

//Connect establishes connection to the database server hosted in "host"
//and returns db handle for the given dbName
func (mdb *MongoDb) Connect(host string, dbName string) DbHandle {
	if mdb == nil {
		log.Fatal("already connected to a db instance\n")
	}

	clientOptions := options.Client().ApplyURI("mongodb://" + host + "/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	mdb.mdbInstance = client.Database(dbName)
	return mdb
}

func (mdb *MongoDb) Create(entityName string) EntityHandle {
	return collectionInfo{collectionHdl: mdb.mdbInstance.Collection(entityName), entityName: entityName}
}

func (mdb *MongoDb) GetEntity(entityName string) EntityHandle {
	return collectionInfo{collectionHdl: mdb.mdbInstance.Collection(entityName), entityName: entityName}
}

func (mdb *MongoDb) Insert(entity EntityHandle, rowEntry Row) error {
	collectionData := entity.(collectionInfo)
	collectionHdl := collectionData.collectionHdl
	_, err := collectionHdl.InsertOne(ctx, rowEntry)
	return err
}

func (mdb *MongoDb) Remove(entity EntityHandle, id ObjectId) error {
	collectionData := entity.(collectionInfo)
	collectionHdl := collectionData.collectionHdl

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	res, err := collectionHdl.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no documents were deleted")
	}

	return nil
}

func (mdb *MongoDb) Get(entity EntityHandle, fn func(func(interface{}) error) (interface{}, error), id ObjectId) (Row, error) {

	var row Row

	collectionData := entity.(collectionInfo)
	collectionHdl := collectionData.collectionHdl

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	cur, err := collectionHdl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if cur.Next(ctx) {
		row, err = fn(cur.Decode)
		if err != nil {
			return nil, err
		}
	}
	cur.Close(ctx)

	return row, nil
}

func filterRows(coll collectionInfo, fn func(func(interface{}) error) (interface{}, error), filter bson.D) ([]Row, error) {
	var rows []Row

	cur, err := coll.collectionHdl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Iterate through the cursor and decode each document one at a time
	for cur.Next(ctx) {

		row, err := fn(cur.Decode)
		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if len(rows) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	fmt.Println()

	return rows, nil
}

func (mdb *MongoDb) GetMatched(entity EntityHandle, fn func(func(interface{}) error) (interface{}, error), match Match) ([]Row, error) {
	var (
		rows   []Row
		err    error
		filter bson.D
	)

	collectionData := entity.(collectionInfo)

	for k, v := range match {
		filter = append(filter, primitive.E{Key: k, Value: v})
	}

	rows, err = filterRows(collectionData, fn, filter)

	return rows, err
}
