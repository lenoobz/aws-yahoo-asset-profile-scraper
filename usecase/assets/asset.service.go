package assets

import (
	"context"

	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-yahoo-asset-profile-scraper/entities"
	"github.com/lenoobz/aws-yahoo-asset-profile-scraper/usecase/checkpoint"
)

// Service exposure
type Service struct {
	assetRepo         Repo
	checkpointService checkpoint.Service
	log               logger.ContextLog
}

// NewService create new service
func NewService(assetRepo Repo, checkpointService checkpoint.Service, log logger.ContextLog) *Service {
	return &Service{
		assetRepo:         assetRepo,
		checkpointService: checkpointService,
		log:               log,
	}
}

// GetAssetsBySource find all assets by source
func (s *Service) GetAssetsBySource(ctx context.Context, source string) ([]*entities.Asset, error) {
	s.log.Info(ctx, "finding all assets by source", "source", source)
	return s.assetRepo.FindAllAssetsBySource(ctx, source)
}

// GetAssetsBySourceFromCheckpoint gets all assets from checkpoint
func (s *Service) GetAssetsBySourceFromCheckpoint(ctx context.Context, source string, pageSize int64) ([]*entities.Asset, error) {
	s.log.Info(ctx, "getting assets from checkpoint")
	numAssets, err := s.assetRepo.CountAssetsBySource(ctx, source)
	if err != nil {
		s.log.Error(ctx, "count assets failed", "error", err)
	}

	checkpoint, err := s.checkpointService.UpdateCheckpoint(ctx, pageSize, numAssets)
	if err != nil {
		s.log.Error(ctx, "find and update checkpoint failed", "error", err)
	}

	if checkpoint == nil {
		s.log.Error(ctx, "checkpoint is nil", "checkpoint", checkpoint)
		return nil, nil
	}

	return s.assetRepo.FindAssetsBySourceFromCheckpoint(ctx, source, checkpoint)
}
