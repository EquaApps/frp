package models

import (
    "time"
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
