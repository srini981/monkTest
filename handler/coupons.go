package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"monk/database"
	"monk/models"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func validateCoupon(coupon models.Coupon) error {
	if time.Now().After(coupon.ExpiresAt) {
		return errors.New("coupon has expired")
	}
	if coupon.Type != "cart-wise" && coupon.Type != "product-wise" && coupon.Type != "bxgy" {
		return errors.New("invalid coupon type")
	}
	return nil
}
func CreateCouponHandler(w http.ResponseWriter, r *http.Request) {
	var coupon models.Coupon
	if err := json.NewDecoder(r.Body).Decode(&coupon); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	coupon.ID = primitive.NewObjectID()

	if err := validateCoupon(coupon); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check for duplicate coupon
	var existingCoupon models.Coupon
	err := database.Collection.FindOne(context.TODO(), bson.M{"type": coupon.Type, "name": coupon.Name}).Decode(&existingCoupon)
	if err != nil && err != mongo.ErrNoDocuments {
		http.Error(w, "Database error while checking for duplicates", http.StatusInternalServerError)
		return
	}

	if err == nil {
		http.Error(w, "Duplicate coupon exists", http.StatusBadRequest)
		return
	}

	// Insert coupon
	_, err = database.Collection.InsertOne(context.TODO(), coupon)

	if err != nil {
		http.Error(w, "Failed to create coupon", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(coupon)
}

func UpdateCouponHandler(w http.ResponseWriter, r *http.Request) {
	var coupon models.Coupon
	if err := json.NewDecoder(r.Body).Decode(&coupon); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid coupon ID", http.StatusBadRequest)
		return
	}

	if err := validateCoupon(coupon); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": coupon,
	}

	_, err = database.Collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, "Failed to update coupon", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(coupon)
}

func GetCouponsHandler(w http.ResponseWriter, r *http.Request) {
	cur, err := database.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "Failed to retrieve coupons", http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.TODO())

	var coupons []models.Coupon
	for cur.Next(context.TODO()) {
		var coupon models.Coupon
		if err := cur.Decode(&coupon); err != nil {

			http.Error(w, "Failed to decode coupon", http.StatusInternalServerError)
			return
		}
		coupons = append(coupons, coupon)
	}
	fmt.Println(coupons)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coupons)
}

func GetCouponByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	fmt.Println(err.Error())
	fmt.Println(id)

	if err != nil {
		http.Error(w, "Invalid coupon ID", http.StatusBadRequest)
		return
	}

	var coupon models.Coupon
	err = database.Collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&coupon)
	fmt.Println(err.Error())

	if err != nil {
		http.Error(w, "Coupon not found", http.StatusNotFound)
		return
	}

	if err := validateCoupon(coupon); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coupon)
}

func DeleteCouponHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid coupon ID", http.StatusBadRequest)
		return
	}

	_, err = database.Collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Failed to delete coupon", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ApplyCouponHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid coupon ID", http.StatusBadRequest)
		return
	}

	var cart models.Cart
	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		http.Error(w, "Invalid cart data", http.StatusBadRequest)
		return
	}

	var coupon models.Coupon
	err = database.Collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&coupon)
	if err != nil {
		http.Error(w, "Coupon not found", http.StatusNotFound)
		return
	}

	if err := validateCoupon(coupon); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Apply the coupon logic
	updatedCart := cart
	var totalDiscount float64

	if coupon.Type == "cart-wise" {
		details := coupon.Details.(map[string]interface{})
		threshold := details["threshold"].(float64)
		discount := details["discount"].(float64)
		total := 0.0
		for _, item := range cart.Items {
			total += item.Price * float64(item.Quantity)
		}
		if total > threshold {
			totalDiscount = total * (discount / 100)
		}
	} else if coupon.Type == "product-wise" {
		details := coupon.Details.(map[string]interface{})
		productID := details["product_id"].(string)
		discount := details["discount"].(float64)
		for i, item := range cart.Items {
			if item.ProductID == productID {
				totalDiscount += item.Price * float64(item.Quantity) * (discount / 100)
				updatedCart.Items[i].Price -= item.Price * (discount / 100)
			}
		}
	} else if coupon.Type == "bxgy" {
		details := coupon.Details.(map[string]interface{})
		buyProducts := details["buy_products"].([]interface{})
		getProducts := details["get_products"].([]interface{})
		repetitionLimit := int(details["repetition_limit"].(float64))
		buyMap := make(map[string]int)
		getMap := make(map[string]int)
		for _, p := range buyProducts {
			prod := p.(map[string]interface{})
			buyMap[prod["product_id"].(string)] = int(prod["quantity"].(float64))
		}
		for _, p := range getProducts {
			prod := p.(map[string]interface{})
			getMap[prod["product_id"].(string)] = int(prod["quantity"].(float64))
		}
		applicableTimes := repetitionLimit
		for productID, requiredQty := range buyMap {
			for _, item := range cart.Items {
				if item.ProductID == productID {
					times := item.Quantity / requiredQty
					if times < applicableTimes {
						applicableTimes = times
					}
				}
			}
		}
		for productID, freeQty := range getMap {
			for i, item := range cart.Items {
				if item.ProductID == productID {
					updatedCart.Items[i].Quantity += applicableTimes * freeQty
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"updated_cart":   updatedCart,
		"total_discount": totalDiscount,
	})
}

func ApplicableCouponsHandler(w http.ResponseWriter, r *http.Request) {
	var cart models.Cart
	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		http.Error(w, "Invalid cart data", http.StatusBadRequest)
		return
	}

	cur, err := database.Collection.Find(context.TODO(), bson.M{"expires_at": bson.M{"$gte": time.Now()}})
	if err != nil {
		http.Error(w, "Failed to fetch coupons", http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.TODO())

	var applicableCoupons []models.ApplicableCoupon

	for cur.Next(context.TODO()) {
		var coupon models.Coupon
		if err := cur.Decode(&coupon); err != nil {
			http.Error(w, "Failed to decode coupon", http.StatusInternalServerError)
			return
		}

		switch coupon.Type {
		case "cart-wise":
			details := coupon.Details.(map[string]interface{})
			threshold := details["threshold"].(float64)
			discount := details["discount"].(float64)
			total := 0.0
			for _, item := range cart.Items {
				total += item.Price * float64(item.Quantity)
			}
			if total > threshold {
				applicableCoupons = append(applicableCoupons, models.ApplicableCoupon{
					CouponID: coupon.ID.Hex(),
					Type:     "cart-wise",
					Discount: total * (discount / 100),
				})
			}
		case "product-wise":
			details := coupon.Details.(map[string]interface{})
			productID := details["product_id"].(string)
			discount := details["discount"].(float64)
			for _, item := range cart.Items {
				if item.ProductID == productID {
					totalDiscount := item.Price * float64(item.Quantity) * (discount / 100)
					applicableCoupons = append(applicableCoupons, models.ApplicableCoupon{
						CouponID: coupon.ID.Hex(),
						Type:     "product-wise",
						Discount: totalDiscount,
					})
					break
				}
			}
		case "bxgy":
			details := coupon.Details.(map[string]interface{})
			buyProducts := details["buy_products"].([]interface{})
			getProducts := details["get_products"].([]interface{})
			repetitionLimit := int(details["repetition_limit"].(float64))

			buyMap := make(map[string]int)
			getMap := make(map[string]int)
			for _, p := range buyProducts {
				prod := p.(map[string]interface{})
				buyMap[prod["product_id"].(string)] = int(prod["quantity"].(float64))
			}
			for _, p := range getProducts {
				prod := p.(map[string]interface{})
				getMap[prod["product_id"].(string)] = int(prod["quantity"].(float64))
			}

			applicableTimes := repetitionLimit
			for productID, requiredQty := range buyMap {
				for _, item := range cart.Items {
					if item.ProductID == productID {
						times := item.Quantity / requiredQty
						if times < applicableTimes {
							applicableTimes = times
						}
					}
				}
			}

			discount := 0.0
			for productID, freeQty := range getMap {
				for _, item := range cart.Items {
					if item.ProductID == productID {
						discount += float64(applicableTimes*freeQty) * item.Price
						break
					}
				}
			}

			if discount > 0 {
				applicableCoupons = append(applicableCoupons, models.ApplicableCoupon{
					CouponID: coupon.ID.Hex(),
					Type:     "bxgy",
					Discount: discount,
				})
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"applicable_coupons": applicableCoupons,
	})
}
