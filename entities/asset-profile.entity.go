package entities

type AssetProfile struct {
	Ticker  string `json:"ticker,omitempty"`
	Sector  string `json:"sector,omitempty"`
	Country string `json:"country,omitempty"`
}
