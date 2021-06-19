package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AssetProfile struct {
	ID         *primitive.ObjectID `bson:"_id,omitempty"`
	IsActive   bool                `bson:"isActive,omitempty"`
	CreatedAt  int64               `bson:"createdAt,omitempty"`
	ModifiedAt int64               `bson:"modifiedAt,omitempty"`
	Enabled    bool                `bson:"enabled,omitempty"`
	Deleted    bool                `bson:"deleted,omitempty"`
	Ticker     string              `bson:"ticker,omitempty"`
	Sector     string              `bson:"sector,omitempty"`
	Country    string              `bson:"country,omitempty"`
}
