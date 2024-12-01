package router

import (
	"monk/handler"

	"github.com/gorilla/mux"
)

func Init(r *mux.Router) {
	r.HandleFunc("/coupons", handler.CreateCouponHandler).Methods("POST")
	r.HandleFunc("/coupons", handler.GetCouponsHandler).Methods("GET")
	r.HandleFunc("/coupons/{id}", handler.GetCouponByIDHandler).Methods("GET")
	r.HandleFunc("/coupons/{id}", handler.UpdateCouponHandler).Methods("PUT")
	r.HandleFunc("/coupons/{id}", handler.DeleteCouponHandler).Methods("DELETE")
	r.HandleFunc("/applicable-coupons", handler.ApplicableCouponsHandler).Methods("POST")
	r.HandleFunc("/apply-coupon/{id}", handler.ApplyCouponHandler).Methods("POST")
}
