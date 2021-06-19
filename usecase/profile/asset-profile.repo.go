package profile

import (
	"context"

	"github.com/hthl85/aws-yahoo-asset-profile-scraper/entities"
)

///////////////////////////////////////////////////////////
// Profile Repository Interface
///////////////////////////////////////////////////////////

// Reader interface
type Reader interface {
}

// Writer interface
type Writer interface {
	UpsertAssetProfile(ctx context.Context, assetProfile *entities.AssetProfile) error
}

// Repo interface
type Repo interface {
	Reader
	Writer
}
