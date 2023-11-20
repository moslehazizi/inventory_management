// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package db

import (
	"time"
)

type Category struct {
	ID           int64  `json:"id"`
	CategoryName string `json:"category_name"`
	SectionName  string `json:"section_name"`
}

type Good struct {
	ID       int64  `json:"id"`
	Category int64  `json:"category"`
	Model    string `json:"model"`
	Unit     int64  `json:"unit"`
	// must be positive and bigger than zero
	Amount    int64     `json:"amount"`
	GoodDesc  string    `json:"good_desc"`
	CreatedAt time.Time `json:"created_at"`
}

type Unit struct {
	ID        int64  `json:"id"`
	UnitName  string `json:"unit_name"`
	UnitValue int64  `json:"unit_value"`
}
