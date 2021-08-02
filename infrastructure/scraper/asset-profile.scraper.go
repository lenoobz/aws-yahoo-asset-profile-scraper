package scraper

import (
	"context"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/google/uuid"
	corid "github.com/lenoobz/aws-lambda-corid"
	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-yahoo-asset-profile-scraper/config"
	"github.com/lenoobz/aws-yahoo-asset-profile-scraper/entities"
	"github.com/lenoobz/aws-yahoo-asset-profile-scraper/usecase/assets"
	"github.com/lenoobz/aws-yahoo-asset-profile-scraper/usecase/profile"
)

// AssetProfileScraper struct
type AssetProfileScraper struct {
	ScrapeAssetProfileJob *colly.Collector
	assetProfileService   *profile.Service
	assetService          *assets.Service
	log                   logger.ContextLog
	errorTickers          []string
	scrapedTickers        []string
}

// NewAssetProfileScraper create new asset profile scraper
func NewAssetProfileScraper(assetService *assets.Service, assetProfileService *profile.Service, log logger.ContextLog) *AssetProfileScraper {
	scrapeAssetProfileJob := newScraperJob()

	return &AssetProfileScraper{
		ScrapeAssetProfileJob: scrapeAssetProfileJob,
		assetProfileService:   assetProfileService,
		assetService:          assetService,
		log:                   log,
	}
}

// newScraperJob creates a new colly collector with some custom configs
func newScraperJob() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains(config.AllowDomain),
		colly.Async(true),
	)

	// Overrides the default timeout (10 seconds) for this collector
	c.SetRequestTimeout(30 * time.Second)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  config.DomainGlob,
		Parallelism: 2,
		RandomDelay: 2 * time.Second,
	})

	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	return c
}

// configJobs configs on error handler and on response handler for scaper jobs
func (s *AssetProfileScraper) configJobs() {
	s.ScrapeAssetProfileJob.OnError(s.errorHandler)
	s.ScrapeAssetProfileJob.OnScraped(s.scrapedHandler)
	s.ScrapeAssetProfileJob.OnHTML("div[data-test=qsp-profile]", s.processAssetProfileResponse)
}

// ScrapeAssetProfilesByTickers scrape asset profiles by tickers
func (s *AssetProfileScraper) ScrapeAssetProfilesByTickers(tickers []string) {
	ctx := context.Background()

	s.configJobs()

	for _, ticker := range tickers {
		reqContext := colly.NewContext()
		reqContext.Put("ticker", ticker)

		url := config.GetAssetProfileByTickerURL(ticker)

		s.log.Info(ctx, "scraping asset profile", "ticker", ticker)
		if err := s.ScrapeAssetProfileJob.Request("GET", url, nil, reqContext, nil); err != nil {
			s.log.Error(ctx, "scraping asset profile", "error", err, "ticker", ticker)
		}
	}

	s.ScrapeAssetProfileJob.Wait()
}

// ScrapeAllAssetProfilesBySource scrape asset profiles by sources
func (s *AssetProfileScraper) ScrapeAllAssetProfilesBySource(source string) {
	ctx := context.Background()

	s.configJobs()

	assets, err := s.assetService.GetAssetsBySource(ctx, source)
	if err != nil {
		s.log.Error(ctx, "scraping asset profile failed", "error", err)
		return
	}

	for _, asset := range assets {
		reqContext := colly.NewContext()
		reqContext.Put("ticker", asset.Ticker)

		url := config.GetAssetProfileByTickerURL(asset.Ticker)

		s.log.Info(ctx, "scraping asset profile", "ticker", asset.Ticker)
		if err := s.ScrapeAssetProfileJob.Request("GET", url, nil, reqContext, nil); err != nil {
			s.log.Error(ctx, "scraping asset profile", "error", err, "ticker", asset.Ticker)
		}
	}

	s.ScrapeAssetProfileJob.Wait()
}

// ScrapeAssetProfilesBySourceFromCheckpoint scrape asset profiles by source from checkpoint
func (s *AssetProfileScraper) ScrapeAssetProfilesBySourceFromCheckpoint(source string, pageSize int64) {
	ctx := context.Background()

	s.configJobs()

	assets, err := s.assetService.GetAssetsBySourceFromCheckpoint(ctx, source, pageSize)
	if err != nil {
		s.log.Error(ctx, "scraping asset profile failed", "error", err)
		return
	}

	for _, asset := range assets {
		reqContext := colly.NewContext()
		reqContext.Put("ticker", asset.Ticker)

		url := config.GetAssetProfileByTickerURL(asset.Ticker)

		s.log.Info(ctx, "scraping asset profile", "ticker", asset.Ticker)
		if err := s.ScrapeAssetProfileJob.Request("GET", url, nil, reqContext, nil); err != nil {
			s.log.Error(ctx, "scraping asset profile failed", "error", err, "ticker", asset.Ticker)
		}
	}

	s.ScrapeAssetProfileJob.Wait()
}

///////////////////////////////////////////////////////////
// Scraper Handler
///////////////////////////////////////////////////////////

// errorHandler generic error handler for all scaper jobs
func (s *AssetProfileScraper) errorHandler(r *colly.Response, err error) {
	ctx := context.Background()
	s.log.Error(ctx, "failed to request url", "url", r.Request.URL, "error", err)
	s.errorTickers = append(s.errorTickers, r.Request.Ctx.Get("ticker"))
}

func (s *AssetProfileScraper) scrapedHandler(r *colly.Response) {
	ctx := context.Background()
	foundSector := r.Ctx.Get("foundSector")
	if foundSector == "" {
		s.log.Error(ctx, "sector not found", "ticker", r.Request.Ctx.Get("ticker"))
		s.errorTickers = append(s.errorTickers, r.Request.Ctx.Get("ticker"))
		return
	}

	foundCountry := r.Ctx.Get("foundCountry")
	if foundCountry == "" {
		s.log.Error(ctx, "country not found", "ticker", r.Request.Ctx.Get("ticker"))
		s.errorTickers = append(s.errorTickers, r.Request.Ctx.Get("ticker"))
		return
	}
}

func (s *AssetProfileScraper) processAssetProfileResponse(e *colly.HTMLElement) {
	// create correlation if for processing fund list
	id, _ := uuid.NewRandom()
	ctx := corid.NewContext(context.Background(), id)

	ticker := e.Request.Ctx.Get("ticker")
	s.log.Info(ctx, "processAssetProfileResponse", "ticker", ticker)

	foundSector := false
	foundCountry := false

	assetProfile := entities.AssetProfile{
		Ticker: ticker,
	}

	e.ForEach("p", func(_ int, paragraph *colly.HTMLElement) {
		if foundCountry {
			return
		}

		var address []string
		paragraph.DOM.Contents().Not("br").Not("a").Each(func(i int, n *goquery.Selection) {
			if goquery.NodeName(n) == "#text" {
				address = append(address, n.Text())
			}
		})

		if len(address) > 0 {
			foundCountry = true
			assetProfile.Country = address[len(address)-1]
		}
	})

	e.ForEach("span", func(_ int, span *colly.HTMLElement) {
		if foundSector {
			return
		}

		if strings.EqualFold(span.Text, "Sector(s)") {
			firstSibling := span.DOM.Siblings().First()
			profileSector := firstSibling.Text()

			if profileSector != "" {
				foundSector = true
				assetProfile.Sector = profileSector
			}
		}
	})

	if foundSector && foundCountry {
		e.Response.Ctx.Put("foundCountry", "true")
		e.Response.Ctx.Put("foundSector", "true")

		if err := s.assetProfileService.AddAssetProfile(ctx, &assetProfile); err != nil {
			s.log.Error(ctx, "add asset profile failed", "error", err, "ticker", assetProfile.Ticker)
			s.errorTickers = append(s.errorTickers, assetProfile.Ticker)
		} else {
			s.scrapedTickers = append(s.scrapedTickers, assetProfile.Ticker)
		}
	}
}

// Close scraper
func (s *AssetProfileScraper) Close() []string {
	s.log.Info(context.Background(), "DONE - SCRAPING ASSET PROFILES", "errorTickers", s.errorTickers)
	return s.scrapedTickers
}
