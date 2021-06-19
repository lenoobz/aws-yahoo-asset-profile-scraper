package assets

import (
	"context"

	"github.com/hthl85/aws-yahoo-asset-profile-scraper/entities"
)

///////////////////////////////////////////////////////////
// Assets Repository Interface
///////////////////////////////////////////////////////////

// Reader interface
type Reader interface {
	FindAssetsBySource(context.Context, string) ([]*entities.Asset, error)
}

// Writer interface
type Writer interface {
}

// Repo interface
type Repo interface {
	Reader
	Writer
}
