package config

import (
	"fmt"
)

// AllowDomain const
const AllowDomain = "ca.finance.yahoo.com"

// DomainGlob const
const DomainGlob = "*yahoo.*"

// GetAssetProfileByTickerURL get sector url
func GetAssetProfileByTickerURL(ticker string) string {
	return fmt.Sprintf("https://ca.finance.yahoo.com/quote/%s/profile?p=%s", ticker, ticker)
}
