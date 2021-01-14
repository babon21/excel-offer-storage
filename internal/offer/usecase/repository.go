package usecase

import "github.com/babon21/excel-offer-storage/internal/offer/domain"

type OfferRepository interface {
	GetList(sellerId string, offerId string, offerName string) ([]domain.Offer, error)
	GetListBySellerId(sellerId string) ([]domain.Offer, error)

	SaveList([]domain.Offer) error
	DeleteList(sellerId string, offerIds []string) error
	UpdateList(offers []domain.Offer) error
}
