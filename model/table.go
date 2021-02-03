package model

import "gorm.io/gorm"

type Table struct {
	gorm.Model
	Number int
	Seats  int
}
