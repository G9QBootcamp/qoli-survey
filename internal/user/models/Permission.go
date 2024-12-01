package models

type Permission struct {
	ID     uint
	Action string `gorm:"unique;not null"`
}
