package models

type Blog struct {
	ID      string `gorm:"primaryKey"`
	Title   string
	Content string
	Author  string
}