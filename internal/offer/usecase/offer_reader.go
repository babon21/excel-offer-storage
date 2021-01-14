package usecase

import "github.com/babon21/excel-offer-storage/internal/offer/domain"

type OfferReader interface {
	ReadAll(filePath string) ([]domain.Offer, uint32, error)
}
