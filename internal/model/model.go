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
	Params    Params `json:"params"`
	UserId    string `json:"user_id"`
	Event     string `json:"event"`
	Timestamp int64  `json:"timestamp"`
	Id        int64  `json:"id"`
}

type Params map[string]any

// // Scan implements the Scanner interface.
func (p *Params) Scan(value any) error {
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
func (p Params) Value() (driver.Value, error) {
	b, err := json.Marshal(&p)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}
