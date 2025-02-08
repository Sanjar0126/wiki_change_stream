package repo

import (
	"context"

	"github.com/Sanjar0126/wiki_change_stream/models"
)

var (
	DiscordUserCollection = "discord_users"
)

type DiscordUserI interface {
	Create(ctx context.Context, req models.DiscordUser) (string, error)
	Update(ctx context.Context, req models.DiscordUser) (string, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*models.DiscordUser, error)
	GetOrCreate(ctx context.Context, id string) (*models.DiscordUser, error)
}
