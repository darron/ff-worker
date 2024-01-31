package d1adaptor

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/darron/ff/core"
	"github.com/google/uuid"

	"github.com/syumai/workers/cloudflare/d1"
	_ "github.com/syumai/workers/cloudflare/d1" // register driver
)

type NewsStoryRepository struct {
	Filename string
}

func (nsr NewsStoryRepository) Connect(dbName string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	c, err := d1.OpenConnector(context.Background(), dbName)
	if err != nil {
		return db, fmt.Errorf("failed to initialize DB: %v", err)
	}
	// use sql.OpenDB instead of sql.Open.
	db = sql.OpenDB(c)
	return db, err
}

func (nsr NewsStoryRepository) Find(id string) (*core.NewsStory, error) {
	ns := core.NewsStory{}
	ctx, cancel := context.WithTimeout(context.Background(), sqliteTimeout)
	defer cancel()
	client, err := nsr.Connect(nsr.Filename)
	if err != nil {
		return &ns, fmt.Errorf("Find/Connect Error: %w", err)
	}
	defer client.Close()
	return nsr.find(ctx, id, client)
}

func (nsr NewsStoryRepository) find(ctx context.Context, id string, client *sql.DB) (*core.NewsStory, error) {
	ns := core.NewsStory{}
	err := client.Get(&ns, "SELECT * from news_stories WHERE id = ?", id)
	if err != nil {
		return &ns, fmt.Errorf("find/Get Error: %w", err)
	}
	return &ns, err
}

func (nsr NewsStoryRepository) Store(ns *core.NewsStory) (string, error) {
	if ns.ID == "" {
		ns.ID = uuid.NewString()
	}
	ctx, cancel := context.WithTimeout(context.Background(), sqliteTimeout)
	defer cancel()
	client, err := nsr.Connect(nsr.Filename)
	if err != nil {
		return "", fmt.Errorf("Store/Connect Error: %w", err)
	}
	defer client.Close()

	return nsr.store(ctx, ns, client)
}

func (nsr NewsStoryRepository) store(ctx context.Context, ns *core.NewsStory, client *sql.DB) (string, error) {
	// Start the transaction.
	tx, err := client.Begin()
	if err != nil {
		return "", fmt.Errorf("store/client.Begin Error: %w", err)
	}

	// Insert/Upsert the NewsStory
	newsStoryQuery := `INSERT INTO news_stories (id, record_id, url, body_text, ai_summary) VALUES (?, ?, ?, ?, ?)
		ON CONFLICT (id)
		DO UPDATE SET url=excluded.url, body_text=excluded.body_text, ai_summary=excluded.ai_summary`
	_, err = tx.Exec(newsStoryQuery, ns.ID, ns.RecordID, ns.URL, ns.BodyText, ns.AISummary)
	if err != nil {
		tx.Rollback() //nolint
		return "", fmt.Errorf("store/NewsStories/tx.Exec Error: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback() //nolint
		return "", fmt.Errorf("store/tx.Commit Error: %w", err)
	}

	return ns.ID, err
}

func (nsr NewsStoryRepository) GetAll() ([]*core.NewsStory, error) {
	var stories []*core.NewsStory

	client, err := nsr.Connect(nsr.Filename)
	if err != nil {
		return stories, fmt.Errorf("GetAll/Connect Error: %w", err)
	}
	defer client.Close()

	err = client.Select(&stories, "SELECT * from news_stories")
	if err != nil {
		return stories, fmt.Errorf("GetAll/Select Error: %w", err)
	}

	return stories, nil
}
