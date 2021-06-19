package assets

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
func NewService(repo Repo, log logger.ContextLog) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

// FindAssetsBySource find all assets by source
func (s *Service) FindAssetsBySource(ctx context.Context, source string) ([]*entities.Asset, error) {
	s.log.Info(ctx, "finding all assets by source", "source", source)
	return s.repo.FindAssetsBySource(ctx, source)
}
