package api

import (
	"github.com/babon21/excel-offer-storage/internal/offer/domain"
	"github.com/babon21/excel-offer-storage/internal/offer/usecase"
)

type DownloadOffersRequest struct {
	SellerId string `json:"seller_id"`
	Url      string `json:"url"`
}

type DownloadOffersResponse struct {
	Statistic usecase.Statistic `json:"statistic"`
}

type AsyncDownloadOffersResponse struct {
	TaskId int64 `json:"task_id"`
}

type GetTaskResponse struct {
	Status    string            `json:"status"`
	Statistic usecase.Statistic `json:"statistic"`
}

type GetListRequest struct {
	SellerId string `json:"seller_id"`
	OfferId  string `json:"offer_id"`
	Name     string `json:"name"`
}

type GetListResponse struct {
	Offers []domain.Offer `json:"offers"`
}
