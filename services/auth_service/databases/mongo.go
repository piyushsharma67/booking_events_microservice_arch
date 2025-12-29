package databases

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/piyushsharma67/events_booking/services/auth_service/models"
	"go.mongodb.org/mongo-driver/bson"
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

// type Database interface {
// 	InsertUser(ctx context.Context, user *models.User) error
// 	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
// }

func (db *MongoDB) InsertUser(ctx context.Context, user *models.UserDocument) error {
	_, err := db.collection.InsertOne(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

func (db *MongoDB) GetUserByEmail(ctx context.Context, email string) (*models.UserDocument, error) {
	var user models.UserDocument

	err := db.collection.FindOne(
		ctx,
		bson.M{
			"email":      email,
			"deleted_at": bson.M{"$exists": false}, // ignore soft-deleted users
		},
	).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // user not found (important distinction)
		}
		return nil, err // real DB error
	}

	return &user, nil
}
