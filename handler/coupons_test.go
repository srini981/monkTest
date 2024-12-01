package handler

import (
	"bytes"
	"encoding/json"
	"monk/database"
	"monk/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCreateCouponHandler(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("valid coupon creation", func(mt *mtest.T) {
		database.InitMongoDB() // Mock initialization
		database.Client = mt.Client
		database.Collection = mt.Coll

		coupon := models.Coupon{
			Type: "cart-wise",
			Details: map[string]interface{}{
				"threshold": 100.0,
				"discount":  10.0,
			},
			Name:      "10% off over $100",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		existingCoupon := bson.D{
			{"type", "cart-wise"},
			{"name", "10% off over $100"},
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.coupons", mtest.FirstBatch), mtest.CreateCursorResponse(1, "coupons", mtest.FirstBatch, existingCoupon)) // Empty response

		payload, _ := json.Marshal(coupon)
		req := httptest.NewRequest(http.MethodPost, "/coupons", bytes.NewReader(payload))
		w := httptest.NewRecorder()

		CreateCouponHandler(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var createdCoupon models.Coupon
		json.Unmarshal(w.Body.Bytes(), &createdCoupon)
		assert.Equal(t, coupon.Type, createdCoupon.Type)
	})
	mt.Run("invalid coupon creation", func(mt *mtest.T) {
		database.InitMongoDB() // Mock initialization
		database.Client = mt.Client
		database.Collection = mt.Coll

		coupon := models.Coupon{
			Type: "cart-wise",
			Details: map[string]interface{}{
				"threshold": 100.0,
				"discount":  10.0,
			},
			Name:      "10% off over $100",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		existingCoupon := bson.D{
			{"type", "cart-wise"},
			{"name", "10% off over $100"},
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.coupons", mtest.FirstBatch, existingCoupon), mtest.CreateCursorResponse(1, "coupons", mtest.FirstBatch, existingCoupon)) // Empty response

		payload, _ := json.Marshal(coupon)
		req := httptest.NewRequest(http.MethodPost, "/coupons", bytes.NewReader(payload))
		w := httptest.NewRecorder()

		CreateCouponHandler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	mt.Run("expired coupon", func(mt *mtest.T) {
		database.InitMongoDB()
		database.Collection = mt.Coll

		coupon := models.Coupon{
			Type: "cart-wise",
			Details: map[string]interface{}{
				"threshold": 100.0,
				"discount":  10.0,
			},
			ExpiresAt: time.Now().Add(-24 * time.Hour),
		}

		payload, _ := json.Marshal(coupon)
		req := httptest.NewRequest(http.MethodPost, "/coupons", bytes.NewReader(payload))
		w := httptest.NewRecorder()

		CreateCouponHandler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "coupon has expired")
	})
}

func TestGetCouponsByIDHandler(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("retrieve  coupons by id ", func(mt *mtest.T) {
		database.InitMongoDB()
		database.Collection = mt.Coll

		existingCoupon := bson.D{
			{"_id", "674c294f3e537627f7301fbc"},
			{"type", "cart-wise"},
			{"name", "10% off over $100"},
		}
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.coupons", mtest.FirstBatch, existingCoupon))

		req := httptest.NewRequest(http.MethodGet, "/coupons/674c294f3e537627f7301fbc", nil)
		w := httptest.NewRecorder()

		GetCouponByIDHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var coupons models.Coupon
		json.Unmarshal(w.Body.Bytes(), &coupons)
		assert.Len(t, coupons, 1)
	})
}

func TestGetCouponsHandler(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("retrieve all coupons", func(mt *mtest.T) {
		database.InitMongoDB()
		database.Collection = mt.Coll

		mockCoupons := []bson.D{{
			{"type", "cart-wise"},
			{"details", bson.M{"threshold": 100, "discount": 10}},
			{"expires_at", time.Now().Add(24 * time.Hour)},
		}, {
			{"type", "cart-wise"},
			{"details", bson.M{"threshold": 100, "discount": 10}},
			{"expires_at", time.Now().Add(24 * time.Hour)},
		},
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.coupons", mtest.FirstBatch, mockCoupons...))

		req := httptest.NewRequest(http.MethodGet, "/coupons", nil)
		w := httptest.NewRecorder()

		GetCouponsHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var coupons []models.Coupon
		json.Unmarshal(w.Body.Bytes(), &coupons)
		assert.Len(t, coupons, len(mockCoupons))
	})
}
