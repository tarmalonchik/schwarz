package service

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"coupon/internal/app/coupon/entity"
	"coupon/internal/app/coupon/repository"
	"coupon/internal/app/coupon/repository/mock"
)

func TestNew(t *testing.T) {
	type args struct {
		repo *mock.Mockrepository
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{"initialize service", args{repo: mock.NewMockrepository(nil)}, Service{repo: mock.NewMockrepository(nil)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.repo); !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ApplyCoupon(t *testing.T) {
	couponString := "TestCoupon"

	type fields struct {
		repo *mock.Mockrepository
	}
	type args struct {
		basket entity.Basket
		code   string
	}

	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockrepository(ctrl)

	repoMock.
		EXPECT().
		FindCouponByCode(couponString).
		DoAndReturn(func(code string) (entity.Coupon, error) {
			return entity.Coupon{
				ID:             uuid.NewString(),
				Code:           couponString,
				Discount:       10,
				MinBasketValue: 80,
			}, nil
		}).
		AnyTimes()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantB   entity.Basket
		wantErr bool
	}{
		{
			"Positive case",
			fields{repoMock},
			args{
				basket: entity.Basket{
					Value:                 100,
					AppliedDiscount:       0,
					ApplicationSuccessful: false,
				},
				code: couponString,
			},
			entity.Basket{
				Value:                 100,
				AppliedDiscount:       10,
				ApplicationSuccessful: true,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				repo: tt.fields.repo,
			}
			gotB, err := s.ApplyCoupon(tt.args.basket, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyCoupon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("ApplyCoupon() gotB = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func TestService_CreateCoupon(t *testing.T) {
	var coupon = entity.Coupon{
		Code:           "Superdiscount",
		Discount:       10,
		MinBasketValue: 55,
	}

	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockrepository(ctrl)

	gomock.InOrder(
		repoMock.
			EXPECT().
			FindCouponByCode(coupon.Code).
			DoAndReturn(func(code string) (entity.Coupon, error) {
				return entity.Coupon{}, repository.ErrCouponNotFound
			}).
			AnyTimes(),
		repoMock.EXPECT().SaveCoupon(
			gomock.Cond(func(in entity.Coupon) bool {
				return in.Discount == coupon.Discount &&
					in.Code == coupon.Code &&
					in.MinBasketValue == coupon.MinBasketValue && in.ID != ""
			})).
			DoAndReturn(func(coupon entity.Coupon) error {
				return nil
			}).AnyTimes(),
	)

	type fields struct {
		repo *mock.Mockrepository
	}

	type args struct {
		discount       int
		code           string
		minBasketValue int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Apply 10%", fields{repoMock}, args{coupon.Discount, coupon.Code, coupon.MinBasketValue}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.fields.repo)
			err := s.CreateCoupon(tt.args.discount, tt.args.code, tt.args.minBasketValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCoupon() gotErr = %v, wantErr %v", err != nil, tt.wantErr)
			}
		})
	}
}
