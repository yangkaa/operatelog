package models

import (
	"time"
)

var remodels []interface{}

func init() {
	remodels = append(remodels,
		&MkyAuditLog{},
	)
}

// GetModels -
func GetModels() []interface{} {
	return remodels
}

//Model model
type Model struct {
	ID        uint      `json:"-"`
	CreatedAt time.Time `json:"createTime" gorm:"column:create_time"`
	UpdatedAt time.Time `json:"updateTime" gorm:"column:update_time"`
}
