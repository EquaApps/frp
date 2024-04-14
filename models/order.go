package models

import (
    "encoding/json"
    "time"
    "github.com/VaalaCat/frp-panel/utils"
    v1 "github.com/fatedier/frp/pkg/config/v1"
    "github.com/samber/lo"
    "gorm.io/gorm"
)

type Order struct {
    *OrderEntity
}

type OrderEntity struct {
    ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
    OrderID   string `json:"order_id" gorm:"uniqueIndex;not null"`
    UserID    int    `json:"user_id" gorm:"not null"`
    Bandwidth int    `json:"bandwidth" gorm:"not null"`
    OperatorID int   `json:"operator_id" gorm:"not null"`
    Comment   string `json:"comment"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (*Order) TableName() string {
    return "orders"
}
