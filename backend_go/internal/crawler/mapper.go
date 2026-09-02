package crawler

import (
    "database/sql"
    "strings"
)

type MappingEngine struct {
    db *sql.DB
}

func NewMappingEngine(db *sql.DB) *MappingEngine {
    return &MappingEngine{db: db}
}

// AutoMap inspects job title & organization to link category_id and trade_id
func (m *MappingEngine) AutoMap(title, organization string) (*int, *int) {
    text := strings.ToLower(title + " " + organization)

    var categoryID *int
    var tradeID *int

    // 1. Try trade keyword matching
    tradeRows, err := m.db.Query("SELECT id, name FROM trades WHERE is_active = true")
    if err == nil {
        defer tradeRows.Close()
        for tradeRows.Next() {
            var id int
            var name string
            if err := tradeRows.Scan(&id, &name); err == nil {
                cleanName := strings.ToLower(name)
                if len(cleanName) > 2 && strings.Contains(text, cleanName) {
                    tradeID = &id
                    break
                }
            }
        }
    }

    // 2. Try category keyword matching
    catRows, err := m.db.Query("SELECT id, name FROM job_categories WHERE is_active = true")
    if err == nil {
        defer catRows.Close()
        for catRows.Next() {
            var id int
            var name string
            if err := catRows.Scan(&id, &name); err == nil {
                cleanName := strings.ToLower(name)
                if len(cleanName) > 2 && strings.Contains(text, cleanName) {
                    categoryID = &id
                    break
                }
            }
        }
    }

    return categoryID, tradeID
}
