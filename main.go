package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"handler-service/config"
	"handler-service/data"
	"handler-service/internal/news"
	"handler-service/openai"
	"handler-service/tgsubscriber"
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

	monitor := &Monitor{handler}
	handler.monitor = monitor

	monitor.initHandling()
	monitor.listenAndServe("0.0.0.0:8080")
}

func (handler *Handler) addNews(news *data.News) {
	fmt.Printf("%++v\n", news)
	newsData, err := json.Marshal(news)
	if err != nil {
		return
	}

	isExisting := handler.isExistingNews(news)
	if isExisting {
		fmt.Printf("This news is already exist: %++v\n", news)
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
	if len(news.Tickers) > 0 {
		news.Predictions, news.Explanations, err = handler.getPredictions(newsData, news.Tickers)
		if err != nil {
			return
		}
	}

	err = handler.repo.Insert(*news)
	if err != nil {
		fmt.Printf("%++v", err)
	}
	fmt.Printf("%++v\n", news)
}

func (handler *Handler) getTickers(newsData []byte) ([]string, error) {
	result, err := openai.ObtainRequest(fmt.Sprintf(GETTICKERS, newsData))
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s\n", result)
	if len(result) == 0 {
		return nil, nil
	}
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
			if len(news.Tickers) <= i {
				break
			}
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

func (handler *Handler) isExistingNews(news *data.News) bool {
	data, err := handler.repo.Get(50, 0)
	if err != nil {
		return false
	}

	for _, n := range data {
		text, err := openai.ObtainRequest(fmt.Sprintf(COMPARENEWS, news, n))
		if err != nil {
			continue
		}
		if text == "ДА" {
			return true
		}
	}

	return false
}

func (handler *Handler) subscribeChannel(link string) error {
	return tgsubscriber.SubscribeChannel(link)
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
