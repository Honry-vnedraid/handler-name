package news

import (
	"database/sql"
	"fmt"
	"handler-service/config"
	"handler-service/data"
	"log"
	"time"

	_ "github.com/lib/pq"
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

func (r *Repository) Insert(title string, text string, t time.Time, source, url string) error {
	_, err := r.db.Exec(`
		INSERT INTO news (title, text, time, source, url)
		VALUES ($1, $2, $3, $4, $5)
	`, title, text, t, source, url)
	return err
}

func (r *Repository) Get(limit int, offset int) ([]data.News, error) {
	query := `
		SELECT title, text, time, source, url 
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
