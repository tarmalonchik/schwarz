package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"coupon/internal/app/coupon/entity"
	"coupon/internal/app/coupon/repository"
	"coupon/internal/app/coupon/service"
	"coupon/internal/pkg/cache"
)

type Config struct {
	Host                     string        `envconfig:"HOST" default:"localhost"`
	Port                     string        `envconfig:"PORT" default:"8080"`
	CouponExpiryDate         time.Time     `envconfig:"PORT" default:"2026-01-19T19:02:25+04:00"`
	GracefulShutdownDuration time.Duration `envconfig:"GRACEFUL_SHUTDOWN_DURATION" default:"5s"`
}

type svc interface {
	ApplyCoupon(basket entity.Basket, code string) (entity.Basket, error)
	CreateCoupon(discount int, code string, minBasketValue int) error
	GetCoupons(codes []string) (coupons []entity.Coupon, err error)
}

type API struct {
	srv  *http.Server
	MUX  *gin.Engine
	conf Config
	svc  svc
}

func New(cfg Config) *API {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	api := &API{
		MUX:  router,
		conf: cfg,
	}
	repo := repository.New(cache.New[string, entity.Coupon]())
	api.svc = service.New(repo)

	api.setServer()
	api.setRoutes()

	return api
}

func (a *API) setServer() {
	a.srv = &http.Server{
		Addr:         fmt.Sprintf(":%s", a.conf.Port),
		Handler:      a.MUX,
		WriteTimeout: time.Second * 15, // to prevent Slow Lori attack we need to set up timeouts
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
}

func (a *API) setRoutes() {
	apiGroup := a.MUX.Group("/api")
	apiGroup.POST("/apply", a.Apply)
	apiGroup.POST("/create", a.Create)
	apiGroup.GET("/coupons", a.Get)
}

func waitForInterruptionSignal(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	cancel()
}

func (a *API) waitForDateInterruption(cancel context.CancelFunc) {
	nowTime := time.Now().UTC()
	duration := a.conf.CouponExpiryDate.Sub(nowTime)
	if duration <= 0 {
		log.Println("Coupon service server alive for a year, closing")
		cancel()
		return
	}
	<-time.After(duration)
	log.Println("Coupon service server alive for a year, closing")
	cancel()
}

func (a *API) Start() {
	ctx, cancel := context.WithCancel(context.Background())

	go waitForInterruptionSignal(cancel)
	go a.waitForDateInterruption(cancel)

	go func() {
		if err := a.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server unexpected error: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Println("Server trying graceful shutdown")
	timeoutContext, timeoutContextCancel := context.WithTimeout(context.Background(), a.conf.GracefulShutdownDuration)
	defer timeoutContextCancel()
	if err := a.srv.Shutdown(timeoutContext); err != nil {
		log.Printf("Server unexpected error during graceful shutdown: %v", err)
	}
	log.Println("Server graceful shutdown done")
}
