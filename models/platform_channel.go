package models

import (
	"time"

	"github.com/imiskolee/optional"
)

type PlatformChannel struct {
	ID        optional.Int64  `json:"id" gorm:"column:id; type:serial;primary_key;"`
	Platform  optional.String `json:"platform" gorm:"not null; type:VARCHAR(255)"`
	Channel   optional.String `json:"Channel" gorm:"not null; type:VARCHAR(255)"`
	CreatedAt time.Time       `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"type:timestamp"`
	Flag      optional.String `json:"flag" gorm:"type:varchar(255)"`
}

func (PlatformChannel) TableName() string {
	return "platform_channels"
}
