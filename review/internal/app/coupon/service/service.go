package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"coupon/internal/app/coupon/entity"
	"coupon/internal/app/coupon/repository"
)

var (
	ErrCouponNotFound                     = fmt.Errorf("coupon not found")
	ErrTryingToApplyCouponToNegativeValue = fmt.Errorf("tried to apply discount to negative value")
	ErrCouponAlreadyExists                = fmt.Errorf("coupon already exists")
	ErrBasketValueIsTooLow                = fmt.Errorf("basket value is too low")
)

type storage interface {
	FindCouponByCode(string) (entity.Coupon, error)
	FindCouponByCodes(codes []string) ([]entity.Coupon, error)
	SaveCoupon(entity.Coupon) error
}

type Service struct {
	repo storage
}

func New(repo storage) *Service {
	return &Service{
		repo: repo,
	}
}

// ApplyCoupon ...
// I am not sure about business logic, how should work applying discount, so I decided to leave this logic as it was
func (s *Service) ApplyCoupon(basket entity.Basket, code string) (entity.Basket, error) {
	coupon, err := s.repo.FindCouponByCode(code)
	if err != nil {
		return entity.Basket{}, ErrCouponNotFound
	}

	if basket.Value < 0 {
		return entity.Basket{}, ErrTryingToApplyCouponToNegativeValue
	}

	if basket.Value < coupon.MinBasketValue {
		return entity.Basket{}, ErrBasketValueIsTooLow
	}

	basket.AppliedDiscount = coupon.Discount
	basket.ApplicationSuccessful = true
	return basket, nil
}

// CreateCoupon ...
// Here is 2 possible scenario, you can replace existing coupon if it exists, or return error. I thought it would be better
// to return error, because it is more obvious behavior, while replacing may be like some magic
func (s *Service) CreateCoupon(discount int, code string, minBasketValue int) error { // we don't need to use here any type
	_, err := s.repo.FindCouponByCode(code)
	if err == nil {
		return ErrCouponAlreadyExists
	}

	if !errors.Is(err, repository.ErrCouponNotFound) {
		return fmt.Errorf("checking coupon exists: %w", err)
	}

	coupon := entity.Coupon{
		ID:             uuid.NewString(),
		Code:           code,
		Discount:       discount,
		MinBasketValue: minBasketValue,
	}

	if err = s.repo.SaveCoupon(coupon); err != nil {
		return fmt.Errorf("saving coupon to repo: %w", err)
	}
	return nil
}

// GetCoupons ...
// It is not a good idea to return err and coupons at the same time
func (s *Service) GetCoupons(codes []string) (coupons []entity.Coupon, err error) {
	if coupons, err = s.repo.FindCouponByCodes(codes); err != nil {
		return nil, fmt.Errorf("getting coupons from repo: %w", err)
	}
	// if len(coupons) != len(codes) {
	// depending on business logic, we can return an error here, or we can continue, or maybe we would like to log this issue
	// }
	return coupons, nil
}
