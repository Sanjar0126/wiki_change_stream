package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Sanjar0126/wiki_change_stream/models"
	"github.com/Sanjar0126/wiki_change_stream/storage/repo"
)

type DiscordUserStorage struct {
	collection *mongo.Collection
}

func NewDiscordUserRepo(db *mongo.Database) repo.DiscordUserI {
	wiki := DiscordUserStorage{
		collection: db.Collection(repo.DiscordUserCollection),
	}

	_, err := wiki.collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "author_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		panic(err)
	}

	return &wiki
}

func (f *DiscordUserStorage) Create(
	ctx context.Context, req models.DiscordUser) (string, error) {
	req.BId = primitive.NewObjectID()

	_, err := f.collection.InsertOne(ctx, req)
	if err != nil {
		return "", err
	}

	return req.BId.Hex(), nil
}

func (f *DiscordUserStorage) Update(
	ctx context.Context, req models.DiscordUser) (string, error) {
	update := bson.M{
		"$set": bson.M{
			"lang": req.Lang,
		},
	}

	filter := bson.M{
		"author_id": bson.M{"$eq": req.AuthorId},
	}

	_, err := f.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return req.BId.Hex(), nil
}

func (f *DiscordUserStorage) Delete(ctx context.Context, id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	result := f.collection.FindOneAndDelete(context.Background(), bson.M{"id": objID})

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (f *DiscordUserStorage) Get(
	ctx context.Context, id string) (*models.DiscordUser, error) {
	var (
		response models.DiscordUser
	)

	if err := f.collection.FindOne(
		ctx,
		bson.M{"author_id": id}).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (f *DiscordUserStorage) GetOrCreate(
	ctx context.Context, id string) (*models.DiscordUser, error) {
	var (
		response models.DiscordUser
	)

	err := f.collection.FindOne(ctx, bson.M{"author_id": id}).Decode(&response)

	if err == mongo.ErrNoDocuments {
		response.BId = primitive.NewObjectID()
		response.AuthorId = id
		response.Lang = "en"

		_, err = f.collection.InsertOne(ctx, response)
	}

	if err != nil {
		return nil, err
	}

	return &response, nil
}
