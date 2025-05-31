package models

type ScoreType struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Type     string `json:"type" gorm:"null"`
	MinScore int    `json:"type" gorm:"null"`
	MaxScore int    `json:"type" gorm:"null"`
}
