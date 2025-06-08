package news

import (
	"database/sql"
	"fmt"
	"handler-service/config"
	"handler-service/data"
	"log"
	"time"

	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func ConnectDB(cfg config.Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("БД недоступна:", err)
	}
	return db
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Insert(news data.News) error {
	_, err := r.db.Exec(`
		INSERT INTO news (title, text, time, source, url, tickers, predictions, explanations)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, news.Title, news.Text, news.Time, news.Source, news.URL,
		pq.Array(news.Tickers), pq.Array(news.Predictions), pq.Array(news.Explanations))
	return err
}

func (r *Repository) Get(limit int, offset int) ([]data.News, error) {
	query := `
		SELECT *
		FROM news
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying news: %v", err)
	}
	defer rows.Close()

	news := make([]data.News, 0)
	for rows.Next() {
		var n data.News
		var t time.Time

		err := rows.Scan(&n.ID, &n.Title, &n.Text, &t, &n.Source, &n.URL,
			pq.Array(&n.Tickers), pq.Array(&n.Predictions), pq.Array(&n.Explanations), &n.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning news: %v", err)
		}

		// Преобразуем time.Time -> string (в формате RFC3339)
		n.Time = t.Format(time.RFC3339)
		news = append(news, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return news, nil
}

func (r *Repository) GetTimeSlice(startDate string, endDate string) ([]data.News, error) {
	query := `
		SELECT *
		FROM news
		WHERE time BETWEEN $1::timestamp AND $2::timestamp;
	`

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error querying news: %v", err)
	}
	defer rows.Close()

	news := make([]data.News, 0)
	for rows.Next() {
		var n data.News
		if err := rows.Scan(&n.Title, &n.Text, &n.Time, &n.Source, &n.URL); err != nil {
			return nil, fmt.Errorf("error scanning news: %v", err)
		}
		news = append(news, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return news, nil
}
