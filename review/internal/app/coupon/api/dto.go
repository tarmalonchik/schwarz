package api

import (
	"coupon/internal/app/coupon/entity"
)

type ApplicationRequest struct {
	Code   string `json:"Code" validate:"required"` // supposed the code should be non-empty
	Basket Basket `json:"Basket"`
}

type Basket struct {
	Value                 int  `json:"Value" validate:"required,gt=0"`           // supposed this value should be positive
	AppliedDiscount       int  `json:"AppliedDiscount" validate:"required,gt=0"` // supposed this value should be positive
	ApplicationSuccessful bool `json:"ApplicationSuccessful"`
}

func (b Basket) ToEntity() entity.Basket {
	return entity.Basket{
		Value:                 b.Value,
		AppliedDiscount:       b.AppliedDiscount,
		ApplicationSuccessful: b.ApplicationSuccessful,
	}
}

type CouponRequest struct {
	Codes []string `json:"Codes" validate:"required"` // no need to make a network request with empty params
}

type Coupon struct {
	Code           string `json:"Code" validate:"required"`                // supposed the code should be non-empty
	Discount       int    `json:"Discount" validate:"required,gt=0"`       // supposed this value should be positive
	MinBasketValue int    `json:"MinBasketValue" validate:"required,gt=0"` // supposed this value should be positive
}

// FillFromEntity using this function to avoid passing sensitive info
func (b Basket) FillFromEntity(basket entity.Basket) Basket {
	return Basket{
		Value:                 basket.Value,
		AppliedDiscount:       basket.AppliedDiscount,
		ApplicationSuccessful: basket.ApplicationSuccessful,
	}
}
