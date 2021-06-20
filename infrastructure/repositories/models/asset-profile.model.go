package models

import (
	"context"
	"time"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-yahoo-asset-profile-scraper/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssetProfileModel struct {
	ID         *primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt  int64               `bson:"createdAt,omitempty"`
	ModifiedAt int64               `bson:"modifiedAt,omitempty"`
	Enabled    bool                `bson:"enabled,omitempty"`
	Deleted    bool                `bson:"deleted,omitempty"`
	Schema     string              `bson:"schema,omitempty"`
	Ticker     string              `bson:"ticker,omitempty"`
	Sector     string              `bson:"sector,omitempty"`
	Country    string              `bson:"country,omitempty"`
}

// NewAssetProfileModel create asset profile model
func NewAssetProfileModel(ctx context.Context, log logger.ContextLog, assetProfile *entities.AssetProfile, schemaVersion string) (*AssetProfileModel, error) {
	return &AssetProfileModel{
		ModifiedAt: time.Now().UTC().Unix(),
		Enabled:    true,
		Deleted:    false,
		Schema:     schemaVersion,
		Ticker:     assetProfile.Ticker,
		Sector:     assetProfile.Sector,
		Country:    assetProfile.Country,
	}, nil
}
