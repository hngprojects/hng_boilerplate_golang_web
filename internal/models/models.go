package models

import (
    "time"
    "github.com/lib/pq"
)

// Testimonial represents a testimonial record
type Testimonial struct {
    ID          int            `json:"id"`
    Author      string         `json:"author"`
    Testimonial string         `json:"testimonial"`
    Comments    pq.StringArray `json:"comments" gorm:"type:text[]"`
    CreatedAt   time.Time      `json:"created_at"`
}

