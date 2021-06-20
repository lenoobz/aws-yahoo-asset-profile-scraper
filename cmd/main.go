package main

import (
	"log"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-yahoo-asset-profile-scraper/config"
	"github.com/hthl85/aws-yahoo-asset-profile-scraper/consts"
	"github.com/hthl85/aws-yahoo-asset-profile-scraper/infrastructure/repositories/repos"
	"github.com/hthl85/aws-yahoo-asset-profile-scraper/infrastructure/scraper"
	"github.com/hthl85/aws-yahoo-asset-profile-scraper/usecase/assets"
	"github.com/hthl85/aws-yahoo-asset-profile-scraper/usecase/profile"
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
	defer assetProfileRepo.Close()

	// create new service
	assetService := assets.NewService(assetRepo, zap)
	profileService := profile.NewService(assetProfileRepo, zap)

	job := scraper.NewAssetProfileScraper(assetService, profileService, zap)
	job.ScrapeAssetsBySource(consts.TIP_RANK_SOURCE)
	defer job.Close()
}
