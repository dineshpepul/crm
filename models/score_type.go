package models

type ScoreType struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Type      string `json:"type" gorm:"null"`
	MinScore  int    `json:"min_score" gorm:"null"`
	MaxScore  int    `json:"max_score" gorm:"null"`
	CompanyId int    `json:"company_id" gorm:"null"`
}
