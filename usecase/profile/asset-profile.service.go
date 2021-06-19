package profile

import (
	"context"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-yahoo-asset-profile-scraper/entities"
)

// Service exposure
type Service struct {
	repo Repo
	log  logger.ContextLog
}

// NewService create new service
func NewService(r Repo, l logger.ContextLog) *Service {
	return &Service{
		repo: r,
		log:  l,
	}
}

// AddAssetProfile add asset profile
func (s *Service) AddAssetProfile(ctx context.Context, assetProfile *entities.AssetProfile) error {
	s.log.Info(ctx, "adding asset profile", "ticker", assetProfile.Ticker)
	return s.repo.UpsertAssetProfile(ctx, assetProfile)
}
