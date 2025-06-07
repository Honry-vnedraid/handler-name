package news

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	"newsapp/config"
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

func (r *Repository) Insert(text string, t time.Time, source, url string) error {
	_, err := r.db.Exec(`
		INSERT INTO news (text, time, source, url)
		VALUES ($1, $2, $3, $4)
	`, text, t, source, url)
	return err
}
