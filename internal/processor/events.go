package processor

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/platoon-cc/platoon-cli/internal/model"
	_ "modernc.org/sqlite"
)

func getDatabaseName(key string) (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return dir, err
	}

	dbFolder := filepath.Join(dir, "platoon")
	os.MkdirAll(dbFolder, 0755)

	return fmt.Sprintf("file:%s/%s.db", dbFolder, key), nil
}

func openDatabase(key string) (*sql.DB, error) {
	dbName, err := getDatabaseName(key)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		return nil, err
	}

	pragma := "PRAGMA foreign_keys = ON;"
	pragma += "PRAGMA journal_mode = WAL;"
	if _, err := db.Exec(pragma); err != nil {
		return nil, err
	}

	migrate := `
CREATE TABLE IF NOT EXISTS events (
  id INTEGER PRIMARY KEY,
  event TEXT,
  user_id TEXT,
  timestamp INTEGER,
  payload JSON
);`

	if _, err := db.Exec(migrate); err != nil {
		return nil, err
	}

	return db, nil
}

type Processor struct {
	db          *sql.DB
	insert_stmt *sql.Stmt
	ingest_stmt *sql.Stmt
	key         string
}

func New(key string) (*Processor, error) {
	p := &Processor{
		key: key,
	}

	db, err := openDatabase(key)
	if err != nil {
		return nil, err
	}

	p.db = db

	p.insert_stmt, err = p.db.Prepare(`INSERT INTO events (id,event,user_id,timestamp,payload) VALUES (?,?,?,?,?)`)
	if err != nil {
		return nil, err
	}

	p.ingest_stmt, err = p.db.Prepare(`INSERT INTO events (event,user_id,timestamp,payload) VALUES (?,?,?,?)`)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Processor) Close() {
	p.db.Close()
}

func (p *Processor) StoreEvents(events []model.Event, idOffset int64) error {
	for _, e := range events {
		eventId := idOffset + e.Id
		t := time.UnixMilli(e.Timestamp).Format("2006/01/02 15:04")
		fmt.Printf("%d %s %s \t%s \t%v\n", eventId, t, e.UserId, e.Event, e.Payload)
		if _, err := p.insert_stmt.Exec(eventId, e.Event, e.UserId, e.Timestamp, e.Payload); err != nil {
			return err
		}
	}
	return nil
}

func (p *Processor) IngestEvent(e model.Event) error {
	t := time.UnixMilli(e.Timestamp).Format("2006/01/02 15:04")
	fmt.Printf("Ingesting: %s (%s) \tuser:%s \tpayload:%v\n", e.Event, t, e.UserId, e.Payload)
	if _, err := p.ingest_stmt.Exec(e.Event, e.UserId, e.Timestamp, e.Payload); err != nil {
		return err
	}
	return nil
}

func (p *Processor) Query(q string) error {
	rows, err := p.db.Query(q)
	if err != nil {
		return err
	}

	for rows.Next() {
		var id string
		var score float32
		if err := rows.Scan(&id, &score); err != nil {
			return err
		}
		fmt.Printf("%s - %f\n", id, score)
	}
	return nil
}

func (p *Processor) GetPeakEventId() (int64, error) {
	row := p.db.QueryRow("select max(id) from events;")
	var res any
	err := row.Scan(&res)
	switch res.(type) {
	case nil:
		return 0, err
	default:
		return res.(int64), err
	}
}
