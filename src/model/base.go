package model

import (
	"time"
)

type BaseEntitySF struct {
	ID int64 `json:"id" gorm:"primary_key;autoIncrement:true;type:bigint"`
	Timestamp
}

type Timestamp struct {
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;type:timestamp;not null;autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp;null;autoUpdateTime"`
}
