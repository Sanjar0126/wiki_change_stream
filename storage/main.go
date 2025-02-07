package storage

import (
	db "go.mongodb.org/mongo-driver/mongo"

	"github.com/Sanjar0126/wiki_change_stream/storage/mongo"
	"github.com/Sanjar0126/wiki_change_stream/storage/repo"
)

type StorageI interface {
	WikiChanges() repo.WikiChangesI
}

type storageMDB struct {
	wikiChangesRepo repo.WikiChangesI
}

func New(db *db.Database) StorageI {
	return &storageMDB{
		wikiChangesRepo: mongo.NewWikiChangesRepo(db),
	}
}

func (s *storageMDB) WikiChanges() repo.WikiChangesI {
	return s.wikiChangesRepo
}
