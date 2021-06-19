package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-yahoo-asset-profile-scraper/config"
	"github.com/hthl85/aws-yahoo-asset-profile-scraper/infrastructure/repositories/repos"
	"github.com/hthl85/aws-yahoo-asset-profile-scraper/infrastructure/scraper"
	"github.com/hthl85/aws-yahoo-asset-profile-scraper/usecase/assets"
	"github.com/hthl85/aws-yahoo-asset-profile-scraper/usecase/profile"
)

func main() {
	lambda.Start(lambdaHandler)
}

func lambdaHandler(ctx context.Context) {
	log.Println("lambda handler is called")

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

	// create new scraper job
	job := scraper.NewAssetProfileScraper(assetService, profileService, zap)
	job.ScrapeAssetsBySource("TIP_RANK")
	defer job.Close()
}
