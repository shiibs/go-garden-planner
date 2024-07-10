package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Garden [][]Plant
func (j Garden) Value() (driver.Value, error) {
    return json.Marshal(j)
}

func (j *Garden) Scan(src interface{}) error {
    return json.Unmarshal(src.([]byte), j)
}

type JSONBDates []time.Time

func (j JSONBDates) Value() (driver.Value, error) {
    return json.Marshal(j)
}

func (j *JSONBDates) Scan(src interface{}) error {
    return json.Unmarshal(src.([]byte), j)
}

func (g Garden) MarshalJSON() ([]byte, error) {
    return json.Marshal([][]Plant(g))
}

func (g *Garden) UnmarshalJSON(data []byte) error {
    return json.Unmarshal(data, (*[][]Plant)(g))
}


func NewGarden(rows, columns int) Garden {
	g := make(Garden, rows)
	for i:= range g {
		g[i] = make([]Plant, columns)
	}

	return g
}


func (g Garden) PlantAt(row, col int, plant Plant) {
	g[row][col] = plant
}