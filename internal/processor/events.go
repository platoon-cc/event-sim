package processor

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/tabwriter"
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
	db    *sql.DB
	key   string
	mutex sync.RWMutex
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

	return p, nil
}

func (p *Processor) Close() {
	p.db.Close()
}

func (p *Processor) StoreEvents(events []model.Event, idOffset int64) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, e := range events {
		eventId := idOffset + e.Id
		t := time.UnixMilli(e.Timestamp).Format("2006/01/02 15:04")
		if _, err := p.db.Exec(`INSERT INTO events (id,event,user_id,timestamp,payload) VALUES (?,?,?,?,?)`, eventId, e.Event, e.UserId, e.Timestamp, e.Payload); err != nil {
			return err
		}
		fmt.Printf("%d %s %s \t%s \t%v\n", eventId, t, e.UserId, e.Event, e.Payload)
	}
	return nil
}

func (p *Processor) IngestEvent(e model.Event) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	t := time.UnixMilli(e.Timestamp).Format("2006/01/02 15:04")
	if _, err := p.db.Exec(`INSERT INTO events (event,user_id,timestamp,payload) VALUES (?,?,?,?)`, e.Event, e.UserId, e.Timestamp, e.Payload); err != nil {
		return err
	}
	fmt.Printf("Ingesting: %s (%s) \tuser:%s \tpayload:%v\n", e.Event, t, e.UserId, e.Payload)

	return nil
}

func (p *Processor) Query2(q string) error {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	rows, err := p.db.Query(q)
	if err != nil {
		return err
	}

	defer rows.Close()

	cols, _ := rows.Columns()

	w := tabwriter.NewWriter(os.Stdout, 0, 2, 1, ' ', 0)
	defer w.Flush()

	sep := []byte("\t")
	newLine := []byte("\n")

	w.Write([]byte(strings.Join(cols, "\t") + "\n"))

	row := make([][]byte, len(cols))
	rowPtr := make([]any, len(cols))
	for i := range row {
		rowPtr[i] = &row[i]
	}

	for rows.Next() {
		_ = rows.Scan(rowPtr...)

		w.Write(bytes.Join(row, sep))
		w.Write(newLine)
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
