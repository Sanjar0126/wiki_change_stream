package repo

import (
	"context"

	"github.com/Sanjar0126/wiki_change_stream/models"
)

var (
	WikiChangesCollection = "wiki_changes"
)

type WikiChangesI interface {
	Create(ctx context.Context, req models.WikiRecentChanges) (string, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*models.WikiRecentChanges, error)
	GetAll(ctx context.Context, offset, limit int64, lang string) (
		[]*models.WikiRecentChanges, int32, error)
}
