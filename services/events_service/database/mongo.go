package database

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/piyushsharma67/events_booking/services/events_service/models"
	"github.com/piyushsharma67/events_booking/services/events_service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	querries   *mongo.Client
	collection *mongo.Collection
}

func NewMongoDb(client *mongo.Client) *MongoDB {
	return &MongoDB{
		querries:   client,
		collection: client.Database(os.Getenv("DATABASE_NAME")).Collection(os.Getenv("DATABASE_COLLECTION")),
	}
}

// making a connection to the mongo db
func ConnectMongo() (*mongo.Client, context.CancelFunc) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Mongo connection failed:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Mongo ping failed:", err)
	}

	log.Println("✅ MongoDB connected")

	return client, cancel
}

//db querries functions

func (db *MongoDB) GenerateEvent(ctx context.Context, event *models.EventDocument) (*models.EventDocument, error) {
	_, err := db.collection.InsertOne(ctx, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (db *MongoDB) DeleteEvent(ctx context.Context, eventId any) error {
	eventIdStr, ok := eventId.(string)
	if !ok {
		return errors.New(utils.EVENT_ID_SHOULD_STRING)
	}
	objID, err := primitive.ObjectIDFromHex(eventIdStr)
	if err != nil {
		return err
	}

	_, err = db.collection.UpdateByID(
		ctx,
		objID,
		bson.M{
			"$set": bson.M{
				"deleted_at": time.Now(),
			},
		},
	)

	return err
}

func (db *MongoDB) UpdateEvent(ctx context.Context, event *models.EventDocument) (*models.EventDocument, error) {
	filter := bson.M{"_id": event.ID}

	update := bson.M{
		"$set": bson.M{
			"title":       event.Title,
			"description": event.Description,
			"updated_at":  time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	var updated models.EventDocument
	err := db.collection.
		FindOneAndUpdate(ctx, filter, update, opts).
		Decode(&updated)

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (db *MongoDB) GetEvent(ctx context.Context, eventId any) (*models.EventDocument, error) {
	eventIdStr, ok := eventId.(string)
	if !ok {
		return nil, errors.New(utils.EVENT_ID_SHOULD_STRING)
	}
	objID, err := primitive.ObjectIDFromHex(eventIdStr)
	if err != nil {
		return nil, err
	}

	var event models.EventDocument
	err = db.collection.
		FindOne(ctx, bson.M{"_id": objID}).
		Decode(&event)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &event, err
}
