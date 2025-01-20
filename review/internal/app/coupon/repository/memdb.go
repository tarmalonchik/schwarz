//go:generate mockgen -source memdb.go -destination=mock/memdb.go -package=mock
package repository

import (
	"fmt"

	"coupon/internal/app/coupon/entity"
)

var (
	ErrCouponNotFound = fmt.Errorf("coupon not found")
)

// DoNotRemove using for mock generation
// nolint
type repository interface {
	FindCouponByCode(code string) (entity.Coupon, error)
	FindCouponByCodes(codes []string) ([]entity.Coupon, error)
	SaveCoupon(coupon entity.Coupon) error
}

type cache interface {
	Save(key string, data entity.Coupon)
	Get(key string) (data entity.Coupon, ok bool)
	GetList(keys []string) []entity.Coupon
}
type Repository struct {
	cache cache
}

func New(cache cache) *Repository {
	return &Repository{
		cache: cache,
	}
}

func (r *Repository) FindCouponByCode(code string) (entity.Coupon, error) { // we don't need to use here pointer, because the struct is too small
	coupon, ok := r.cache.Get(code)
	if !ok {
		return entity.Coupon{}, ErrCouponNotFound
	}
	return coupon, nil
}

// FindCouponByCodes ...
// Usually we are using network to access storage, and it is quite slow to make a lot of network request, so it is better
// to get from storage a bucket of codes in one time.
// It is also possible to set some limitations for a list of codes, because network can also have some limitation for request size,
// but not in this case
func (r *Repository) FindCouponByCodes(codes []string) ([]entity.Coupon, error) {
	return r.cache.GetList(codes), nil
}

func (r *Repository) SaveCoupon(coupon entity.Coupon) error {
	r.cache.Save(coupon.Code, coupon)
	return nil
}
