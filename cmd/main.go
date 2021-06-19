package main

import (
	"log"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-yahoo-asset-profile-scraper/config"
	"github.com/hthl85/aws-yahoo-asset-profile-scraper/infrastructure/repositories/repos"
	"github.com/hthl85/aws-yahoo-asset-profile-scraper/infrastructure/scraper"
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

	// create new service
	profileService := profile.NewService(assetProfileRepo, zap)

	ts := scraper.NewAssetProfileScraper(profileService, zap)
	ts.ScrapeAssetProfilesByTickers([]string{"FAP.TO", "TD.TO"})
	// ts.ScrapeAssetsPriceFromCheckpoint(consts.PAGE_SIZE)
	// defer ts.Close()
}
