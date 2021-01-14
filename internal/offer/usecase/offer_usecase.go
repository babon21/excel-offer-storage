package usecase

import (
	"github.com/babon21/excel-offer-storage/internal/offer/domain"
	"github.com/babon21/excel-offer-storage/internal/offer/reader"
)

type OfferUseCase interface {
	Store(sellerId string, url string) (Statistic, error)
	GetList(sellerId string, offerId string, offerName string) ([]domain.Offer, error)
}

type offerUseCase struct {
	offerRepository OfferRepository
	//offerReader     OfferReader
	offerGateway   OfferGateway
	excelFilesPath string
}

func (useCase *offerUseCase) Store(sellerId string, url string) (Statistic, error) {
	offers, errCount, err := useCase.getOffers(sellerId, url)
	if err != nil {
		return Statistic{}, err
	}

	statistic := Statistic{ErrCount: errCount}

	dbOffers, err := useCase.offerRepository.GetListBySellerId(sellerId)
	if err != nil {
		return statistic, err
	}
	_ = dbOffers

	offersToSave := make([]domain.Offer, 0, 1)
	_ = offersToSave
	offersToDelete := make([]string, 0, 1)
	_ = offersToDelete
	offersToUpdate := make([]domain.Offer, 0, 1)
	_ = offersToUpdate

	offerMap := createOfferMap(offers)
	dbOfferMap := createOfferMap(dbOffers)

	for k, v := range offerMap {
		if _, found := dbOfferMap[k]; !found {
			offersToSave = append(offersToSave, v)
			continue
		}

		if !v.Available {
			offersToDelete = append(offersToDelete, k)
			continue
		}

		offersToUpdate = append(offersToUpdate, v)
	}

	if err := useCase.offerRepository.SaveList(offersToSave); err != nil {
		return statistic, err
	}

	if err = useCase.offerRepository.DeleteList(sellerId, offersToDelete); err != nil {
		return statistic, err
	}

	if err = useCase.offerRepository.UpdateList(offersToUpdate); err != nil {
		return statistic, nil
	}

	statistic.CreatedCount = uint32(len(offersToSave))
	statistic.DeletedCount = uint32(len(offersToDelete))
	statistic.UpdatedCount = uint32(len(offersToUpdate))

	return statistic, nil
}

func createOfferMap(offers []domain.Offer) map[string]domain.Offer {
	var offerMap = make(map[string]domain.Offer, len(offers))
	for _, offer := range offers {
		offerMap[offer.OfferId] = offer
	}
	return offerMap
}

func (useCase *offerUseCase) GetList(sellerId string, offerId string, offerName string) ([]domain.Offer, error) {
	return useCase.offerRepository.GetList(sellerId, offerId, offerName)
}

func (useCase *offerUseCase) getOffers(sellerId string, offerPath string) ([]domain.Offer, uint32, error) {
	offerFilename, err := useCase.offerGateway.DownloadOffers(offerPath)
	defer useCase.offerGateway.DeleteOffers(offerFilename)
	if err != nil {
		return nil, 0, err
	}

	offers, errCount, err := reader.ReadAll(offerFilename)
	if err != nil {
		return nil, errCount, err
	}

	for i := range offers {
		offers[i].SellerId = sellerId
	}

	return offers, errCount, nil
}

func NewOfferUseCase(offerRepository OfferRepository, offerGateway OfferGateway, excelFilesPath string) OfferUseCase {
	return &offerUseCase{
		offerRepository: offerRepository,
		//offerReader:     offerReader,
		offerGateway:   offerGateway,
		excelFilesPath: excelFilesPath,
	}
}
