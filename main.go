package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"handler-service/data"
	"handler-service/internal/news"
	"handler-service/openai"
	"strings"
	"time"
)

type Handler struct {
	monitor *Monitor
	db      *sql.DB
	repo    *news.Repository
}

func main() {
	handler := &Handler{}
	// _ = godotenv.Load()

	// cfg := config.LoadConfig()

	// db := news.ConnectDB(cfg)
	// defer db.Close()

	// repo := news.NewRepository(db)

	// handler := &Handler{
	// 	db:   db,
	// 	repo: repo,
	// }

	// jsons := []string{
	// 	`{
	// 		"title": "1",
	// 		"text": "SpaceX успешно провела испытание системы аварийного спасения экипажа...",
	// 		"time": "2025-06-07T06:50:00Z",
	// 		"source": "BBC Science",
	// 		"url": "https://bbc.com/news/science/spacex-starship-crew-test"
	// 	}`,
	// 	`{
	// 		"title": "2",
	// 		"text": "Ракета Starship от SpaceX прошла ключевой тест...",
	// 		"time": "2025-06-07T07:05:00Z",
	// 		"source": "РИА Новости",
	// 		"url": "https://ria.ru/20250607/spacex-starship-test-1850000000.html"
	// 	}`,
	// }

	// for _, js := range jsons {
	// 	var n data.News
	// 	_ = json.Unmarshal([]byte(js), &n)
	// 	parsedTime, _ := time.Parse(time.RFC3339, n.Time)
	// 	err := repo.Insert(n.Title, n.Text, parsedTime, n.Source, n.URL)
	// 	if err != nil {
	// 		log.Println("❌ Ошибка вставки:", err)
	// 	}
	// }

	// log.Println("✅ Новости добавлены")

	monitor := &Monitor{handler}
	handler.monitor = monitor

	monitor.initHandling()
	monitor.listenAndServe("127.0.0.1:8080")
}

func (handler *Handler) addNews(news *data.News) {
	newsData, err := json.Marshal(news)
	if err != nil {
		return
	}

	// if news.Title == "" {
	// 	news.Title, err = openai.ObtainRequest(fmt.Sprintf(GENERATETITLEPROMPT, newsData))
	// 	if err != nil {
	// 		return
	// 	}
	// }

	news.Tickers, err = handler.getTickers(newsData)
	if err != nil {
		return
	}
	news.Predictions, news.Explanations, err = handler.getPredictions(newsData, news.Tickers)
	if err != nil {
		return
	}

	// time, err := parseDateTime(news.Time)
	// if err != nil {
	// 	return
	// }

	// handler.repo.Insert(news.Title, news.Text, time, news.Source, news.URL)
	fmt.Printf("%++v\n", news)
}

func (handler *Handler) getTickers(newsData []byte) ([]string, error) {
	result, err := openai.ObtainRequest(fmt.Sprintf(GETTICKERS, newsData))
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s\n", result)
	return strings.Split(result, ", "), nil
}

func (handler *Handler) getPredictions(newsData []byte, tickers []string) ([]string, []string, error) {
	result, err := openai.ObtainRequest(fmt.Sprintf(GETPREDICTIONS, newsData, tickers))
	if err != nil {
		return nil, nil, err
	}
	result = strings.TrimSpace(result)
	data := strings.Split(result, "\n")
	preds := strings.TrimSpace(data[0])
	return strings.Split(preds, ", "), data[1:], nil
}

func (handler *Handler) getNews(offset int, limit int) ([]data.News, error) {
	data, err := handler.repo.Get(limit, offset)
	return data, err
}

func parseDateTime(datetimeStr string) (time.Time, error) {
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
		time.RFC3339,
		"Jan 2, 2006 at 3:04pm (MST)",
	}

	var lastErr error

	for _, layout := range layouts {
		t, err := time.Parse(layout, datetimeStr)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}

	return time.Time{}, fmt.Errorf("failed to parse time string '%s': %v", datetimeStr, lastErr)
}
