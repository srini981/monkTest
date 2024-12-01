package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coupon struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type      string             `json:"type"`
	Details   interface{}        `json:"details"`
	Name      string             `json:"name"`
	ExpiresAt time.Time          `json:"expires_at"`
}

type Cart struct {
	Items []CartItem `json:"items"`
}

type CartItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type ApplicableCoupon struct {
	CouponID string  `json:"coupon_id"`
	Type     string  `json:"type"`
	Discount float64 `json:"discount"`
}
