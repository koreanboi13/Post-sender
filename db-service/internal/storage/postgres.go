package storage

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db   *sql.DB
	psql sq.StatementBuilderType
}

func NewPostgres(connectionString string) (*Postgres, error) {
	log.Println(connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	pg := &Postgres{
		db:   db,
		psql: psql,
	}

	if err := pg.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return pg, nil
}

func (p *Postgres) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS chat_entries (
		id VARCHAR(255) PRIMARY KEY,
		messenger VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);
	
	-- Создаем индекс по полю messenger для быстрого поиска
	CREATE INDEX IF NOT EXISTS idx_chat_entries_messenger ON chat_entries(messenger);
	`

	_, err := p.db.Exec(query)
	return err
}

func (p *Postgres) Save(chatID string, messengerType MessengerType) error {
	query := p.psql.Insert("chat_entries").
		Columns("id", "messenger", "created_at").
		Values(chatID, messengerType, sq.Expr("NOW()")).
		Suffix("ON CONFLICT (id) DO UPDATE SET messenger = EXCLUDED.messenger")

	_, err := query.RunWith(p.db).Exec()
	if err != nil {
		return fmt.Errorf("failed to save chat entry: %w", err)
	}

	return nil
}

func (p *Postgres) Delete(chatID string) error {
	query := p.psql.Delete("chat_entries").
		Where(sq.Eq{"id": chatID})

	_, err := query.RunWith(p.db).Exec()
	if err != nil {
		return fmt.Errorf("failed to delete chat entry: %w", err)
	}

	return nil
}

func (p *Postgres) Exists(chatID string) (bool, error) {
	query := p.psql.Select("1").
		From("chat_entries").
		Where(sq.Eq{"id": chatID}).
		Limit(1)

	var exists int
	err := query.RunWith(p.db).QueryRow().Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if chat entry exists: %w", err)
	}

	return true, nil
}

func (p *Postgres) GetAllByMessenger(messengerType MessengerType) ([]int, error) {
	query := p.psql.Select("id").
		From("chat_entries").
		Where(sq.Eq{"messenger": messengerType}).
		OrderBy("created_at DESC")

	rows, err := query.RunWith(p.db).Query()
	if err != nil {
		return nil, fmt.Errorf("failed to query chat entries by messenger: %w", err)
	}
	defer rows.Close()

	var chatIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan chat ID: %w", err)
		}
		chatIDs = append(chatIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over chat IDs: %w", err)
	}
	var Ids []int
	for _, id := range chatIDs {
		temp, err := strconv.Atoi(id)
		if err != nil {
			continue
		}
		Ids = append(Ids, temp)
	}
	return Ids, nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}
