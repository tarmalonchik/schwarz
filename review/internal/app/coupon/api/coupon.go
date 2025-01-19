package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"coupon/internal/app/coupon/service"
)

func (a *API) Apply(c *gin.Context) {
	apiReq := ApplicationRequest{}

	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.JSON(400, "invalid request")
		return
	}

	if err := validator.New().Struct(apiReq); err != nil {
		c.JSON(400, err.Error())
		return
	}

	serviceResponse, err := a.svc.ApplyCoupon(apiReq.Basket.ToEntity(), apiReq.Code)
	if err != nil {
		if errors.Is(err, service.ErrCouponNotFound) {
			c.JSON(404, "coupon not found")
			return
		}
		if errors.Is(err, service.ErrTryingToApplyCouponToNegativeValue) {
			c.JSON(400, "invalid request")
			return
		}
		if errors.Is(err, service.ErrBasketValueIsTooLow) {
			c.JSON(400, "invalid request: basket value is too low")
			return
		}
		log.Printf("service error: %v", err)
		c.JSON(500, "service error")
		return
	}

	var basket Basket

	basket = basket.FillFromEntity(serviceResponse)

	c.JSON(http.StatusOK, basket)
}

func (a *API) Create(c *gin.Context) {
	apiReq := Coupon{}

	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.JSON(400, "invalid request")
		return
	}

	if err := validator.New().Struct(apiReq); err != nil {
		c.JSON(400, err.Error())
		return
	}

	err := a.svc.CreateCoupon(apiReq.Discount, apiReq.Code, apiReq.MinBasketValue)
	if err != nil {
		if errors.Is(err, service.ErrCouponAlreadyExists) {
			c.JSON(409, "coupon already exists")
			return
		}
		log.Printf("service error: %v", err)
		c.JSON(500, "service error")
		return
	}
	c.Status(http.StatusOK)
}

// Get ...
// Usually Get requests are not taking params from body, usually they take arguments from query param, but in this case
// I decided not to change api, because this can affect clients
func (a *API) Get(c *gin.Context) {
	apiReq := CouponRequest{}

	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.JSON(400, "invalid request")
		return
	}

	if err := validator.New().Struct(apiReq); err != nil {
		c.JSON(400, err.Error())
		return
	}

	coupons, err := a.svc.GetCoupons(apiReq.Codes)
	if err != nil {
		log.Printf("service error: %v", err)
		c.JSON(500, "service error")
		return
	}

	c.JSON(http.StatusOK, coupons)
}
