package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"handler-service/config"
	"handler-service/data"
	"handler-service/internal/news"
	"handler-service/openai"
	"log"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Handler struct {
	monitor *Monitor
	db      *sql.DB
	repo    *news.Repository
}

func main() {
	_ = godotenv.Load()

	cfg := config.LoadConfig()

	db := news.ConnectDB(cfg)
	defer db.Close()

	repo := news.NewRepository(db)

	handler := &Handler{
		db:   db,
		repo: repo,
	}

	jsons := []string{
		`{
			"title": "1",
			"text": "SpaceX успешно провела испытание системы аварийного спасения экипажа...",
			"time": "2025-06-07T06:50:00Z",
			"source": "BBC Science",
			"url": "https://bbc.com/news/science/spacex-starship-crew-test"
		}`,
		`{
			"title": "2",
			"text": "Ракета Starship от SpaceX прошла ключевой тест...",
			"time": "2025-06-07T07:05:00Z",
			"source": "РИА Новости",
			"url": "https://ria.ru/20250607/spacex-starship-test-1850000000.html"
		}`,
	}

	for _, js := range jsons {
		var n data.News
		_ = json.Unmarshal([]byte(js), &n)
		err := repo.Insert(n)
		if err != nil {
			log.Println("❌ Ошибка вставки:", err)
		}
	}

	log.Println("✅ Новости добавлены")

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

	if news.Title == "" {
		news.Title, err = openai.ObtainRequest(fmt.Sprintf(GENERATETITLEPROMPT, newsData))
		if err != nil {
			return
		}
	}

	news.Tickers, err = handler.getTickers(newsData)
	if err != nil {
		return
	}
	news.Predictions, news.Explanations, err = handler.getPredictions(newsData, news.Tickers)
	if err != nil {
		return
	}

	handler.repo.Insert(*news)
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

func (handler *Handler) getSummary(startDate string, endDate string) (*data.Summary, error) {
	datas, err := handler.repo.GetTimeSlice(startDate, endDate)
	if err != nil {
		return nil, err
	}

	changes := make(map[string]int, 0)
	for _, news := range datas {
		for i, item := range news.Predictions {
			name := news.Tickers[i]
			val, err := strconv.Atoi(item)
			if err != nil {
				continue
			}
			changes[name] += val
		}

	}

	text, err := openai.ObtainRequest(fmt.Sprintf(GETSUMMARY, datas))
	if err != nil {
		return nil, err
	}

	tickers := make([]string, 0)
	predictions := make([]int, 0)
	for key, value := range changes {
		tickers = append(tickers, key)
		if value > 100 {
			value = 100
		} else if value < -100 {
			value = -100
		}
		predictions = append(predictions, value)
	}

	result := makeSummary(text, tickers, predictions)

	return result, nil
}

func makeSummary(
	text string,
	tickers []string,
	predictions []int,
) *data.Summary {
	return &data.Summary{
		Text:        text,
		Tickers:     tickers,
		Predictions: predictions,
	}
}
