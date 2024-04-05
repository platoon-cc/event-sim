package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Team struct {
	Id   string
	Name string
}

type Project struct {
	Id   string
	Name string
}

type Event struct {
	Payload   Payload `json:"payload"`
	UserId    string  `json:"user_id"`
	Event     string  `json:"event"`
	Timestamp int64   `json:"timestamp"`
	Id        int64   `json:"id"`
}

type Payload map[string]any

// // Scan implements the Scanner interface.
func (p *Payload) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		json.Unmarshal(v, p)
		return nil
	case string:
		json.Unmarshal([]byte(v), p)
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

// // Value implements the driver Valuer interface.
func (p Payload) Value() (driver.Value, error) {
	b, err := json.Marshal(&p)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}
