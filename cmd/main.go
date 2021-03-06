package main

import (
	"log"

	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-yahoo-asset-profile-scraper/config"
	"github.com/lenoobz/aws-yahoo-asset-profile-scraper/consts"
	"github.com/lenoobz/aws-yahoo-asset-profile-scraper/infrastructure/repositories/repos"
	"github.com/lenoobz/aws-yahoo-asset-profile-scraper/infrastructure/scraper"
	"github.com/lenoobz/aws-yahoo-asset-profile-scraper/usecase/assets"
	"github.com/lenoobz/aws-yahoo-asset-profile-scraper/usecase/checkpoint"
	"github.com/lenoobz/aws-yahoo-asset-profile-scraper/usecase/profile"
)

func main() {
	appConf := config.AppConf

	// create new logger
	zap, err := logger.NewZapLogger()
	if err != nil {
		log.Fatal("create app logger failed")
	}
	defer zap.Close()

	// create new repository
	assetProfileRepo, err := repos.NewAssetProfileMongo(nil, zap, &appConf.Mongo)
	if err != nil {
		log.Fatal("create asset profile mongo failed")
	}
	defer assetProfileRepo.Close()

	// create new repository
	assetRepo, err := repos.NewAssetMongo(nil, zap, &appConf.Mongo)
	if err != nil {
		log.Fatal("create asset mongo failed")
	}
	defer assetRepo.Close()

	// create new repository
	checkpointRepo, err := repos.NewCheckpointMongo(nil, zap, &appConf.Mongo)
	if err != nil {
		log.Fatal("create checkpoint mongo failed")
	}
	defer checkpointRepo.Close()

	// create new service
	checkpointService := checkpoint.NewService(checkpointRepo, zap)
	assetService := assets.NewService(assetRepo, *checkpointService, zap)
	profileService := profile.NewService(assetProfileRepo, zap)

	job := scraper.NewAssetProfileScraper(assetService, profileService, zap)
	// job.ScrapeAllAssetProfilesBySource(consts.TIP_RANK_SOURCE)
	job.ScrapeAssetProfilesBySourceFromCheckpoint(consts.TIP_RANK_SOURCE, consts.PAGE_SIZE)
	defer job.Close()
}
