package data

type Summary struct {
	Text        string   `json:"text"`
	Tickers     []string `json:"tickers"`
	Predictions []string `json:"predictions"`
}
