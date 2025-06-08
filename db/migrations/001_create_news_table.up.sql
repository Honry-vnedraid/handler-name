CREATE TABLE IF NOT EXISTS news (
	id SERIAL PRIMARY KEY,
	title TEXT,
	text TEXT,
	source TEXT,
	url TEXT,
	tickers TEXT[],
	predictions TEXT[],
	explanations TEXT[],
	time TIMESTAMP DEFAULT now(),
    created_at TIMESTAMP DEFAULT now()
);
