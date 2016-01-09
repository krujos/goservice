package main

import "time"

type Widget struct {
	Description string  `json:"description"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	Created   time.Time `json:"due"`
}

type Widgets []Widget
