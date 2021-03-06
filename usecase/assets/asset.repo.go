package assets

import (
	"context"

	"github.com/lenoobz/aws-yahoo-asset-profile-scraper/entities"
)

///////////////////////////////////////////////////////////
// Assets Repository Interface
///////////////////////////////////////////////////////////

// Reader interface
type Reader interface {
	CountAssetsBySource(context.Context, string) (int64, error)
	FindAllAssetsBySource(context.Context, string) ([]*entities.Asset, error)
	FindAssetsBySourceFromCheckpoint(context.Context, string, *entities.Checkpoint) ([]*entities.Asset, error)
}

// Writer interface
type Writer interface {
}

// Repo interface
type Repo interface {
	Reader
	Writer
}
