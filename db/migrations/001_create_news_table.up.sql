CREATE TABLE IF NOT EXISTS news (
	id SERIAL PRIMARY KEY,
	title TEXT,
	text TEXT,
	time TIMESTAMP,
	source TEXT,
	url TEXT,
	tickers TEXT[],
	predictions TEXT[],
	explanations TEXT[],
    created_at TIMESTAMP DEFAULT now()
);
