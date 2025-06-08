package data

type News struct {
	Title        string   `json:"title"`
	Text         string   `json:"text"`
	Time         string   `json:"time"`
	Source       string   `json:"source"`
	URL          string   `json:"url"`
	Tickers      []string `json:"tickers"`
	Predictions  []string `json:"predictions"`
	Explanations []string `json:"explanations"`
}
