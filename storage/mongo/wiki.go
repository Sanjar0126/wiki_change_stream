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

type wikiChangesStorage struct {
	collection *mongo.Collection
}

func NewWikiChangesRepo(db *mongo.Database) repo.WikiChangesI {
	wiki := wikiChangesStorage{
		collection: db.Collection(repo.WikiChangesCollection),
	}

	_, err := wiki.collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "meta.id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		panic(err)
	}

	return &wiki
}

func (f *wikiChangesStorage) Create(
	ctx context.Context, req models.WikiRecentChanges) (string, error) {
	req.BId = primitive.NewObjectID()

	_, err := f.collection.InsertOne(ctx, req)
	if err != nil {
		return "", err
	}

	return req.BId.Hex(), nil
}

func (f *wikiChangesStorage) Delete(ctx context.Context, id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	result := f.collection.FindOneAndDelete(context.Background(), bson.M{"id": objID})

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (f *wikiChangesStorage) Get(
	ctx context.Context, id string) (*models.WikiRecentChanges, error) {
	var (
		response models.WikiRecentChanges
	)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if err = f.collection.FindOne(
		ctx,
		bson.M{"id": objectID}).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (f *wikiChangesStorage) GetAll(ctx context.Context, page, limit int64, lang string) (
	[]*models.WikiRecentChanges, int32, error) {
	var (
		response []*models.WikiRecentChanges
	)

	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSkip((page - 1) * limit)

	filtering := bson.M{"server_prefix": lang}

	count, err := f.collection.CountDocuments(context.Background(), filtering)
	if err != nil {
		return response, 0, err
	}

	rows, err := f.collection.Find(
		context.Background(),
		filtering,
		opts,
	)
	if err != nil {
		return response, 0, err
	}

	if err := rows.All(context.Background(), &response); err != nil {
		return response, 0, err
	}

	return response, int32(count), nil
}
