CREATE TABLE IF NOT EXISTS urls (
    id BIGSERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(20) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP,
    click_count BIGINT NOT NULL DEFAULT 0
);

-- Индекс для быстрого поиска по short_code
CREATE INDEX idx_urls_short_code ON urls(short_code);

-- Индекс для поиска истекших ссылок
CREATE INDEX idx_urls_expires_at ON urls(expires_at) WHERE expires_at IS NOT NULL;

-- Индекс для сортировки по дате создания
CREATE INDEX idx_urls_created_at ON urls(created_at DESC);

-- Комментарии к таблице и колонкам
COMMENT ON TABLE urls IS 'Таблица для хранения сокращенных URL';
COMMENT ON COLUMN urls.id IS 'Уникальный идентификатор';
COMMENT ON COLUMN urls.original_url IS 'Исходный полный URL';
COMMENT ON COLUMN urls.short_code IS 'Короткий уникальный код';
COMMENT ON COLUMN urls.created_at IS 'Дата и время создания';
COMMENT ON COLUMN urls.expires_at IS 'Дата и время истечения (NULL = бессрочно)';
COMMENT ON COLUMN urls.click_count IS 'Количество переходов по ссылке';