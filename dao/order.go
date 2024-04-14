package dao

import (
    "github.com/EquaApps/frp/models"
    "time"
	"github.com/google/uuid"
)

func CreateOrder(userID, bandwidth, operatorID int, comment string) (*models.OrderEntity, error) {
    db := models.GetDBManager().GetDefaultDB()

    now := time.Now()
    order := &models.OrderEntity{
        OrderID:    generateOrderID(),
        UserID:     userID,
        Bandwidth:  bandwidth,
        OperatorID: operatorID,
        Comment:    comment,
        CreatedAt:  now,
        UpdatedAt:  now,
    }

    err := db.Create(order).Error
    if err != nil {
        return nil, err
    }

    return order, nil
}

func generateOrderID() string {
    // 生成订单ID的逻辑
    return "ORD-" + uuid.New().String()
}