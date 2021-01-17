package main

import (
	"fmt"
	"github.com/babon21/excel-offer-storage/internal/config"
	"github.com/babon21/excel-offer-storage/internal/http/middleware"
	offerHttp "github.com/babon21/excel-offer-storage/internal/offer/delivery/http"
	"github.com/babon21/excel-offer-storage/internal/offer/gateway"
	"github.com/babon21/excel-offer-storage/internal/offer/reader"
	"github.com/babon21/excel-offer-storage/internal/offer/repository/postgres"
	"github.com/babon21/excel-offer-storage/internal/offer/store"
	"github.com/babon21/excel-offer-storage/internal/offer/usecase"
	asyncUsecase "github.com/babon21/excel-offer-storage/internal/offer/usecase/async"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

func main() {
	conf := config.Init()

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", conf.Database.Username,
		conf.Database.Password, conf.Database.Host, conf.Database.Port, conf.Database.DbName)
	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.AccessLogMiddleware)
	offerRepo := postgres.NewPostgresOfferRepository(db)
	gw := gateway.NewOfferGateway(".")
	offerUseCase := usecase.NewOfferUseCase(offerRepo, gw, reader.NewExcelOfferReader())
	redisStore, err := store.NewRedisStore(conf.Cache.Host, conf.Cache.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while connecting to redis")
	}

	asyncUseCase := asyncUsecase.NewAsyncOfferUseCase(offerUseCase, redisStore)
	offerHttp.NewOfferHandler(e, offerUseCase, asyncUseCase)
	log.Fatal().Msg(e.Start(":" + conf.Server.Port).Error())
}
