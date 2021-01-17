package usecase_test

import (
	"errors"
	"github.com/babon21/excel-offer-storage/internal/offer/domain"
	"github.com/babon21/excel-offer-storage/internal/offer/domain/mocks"
	"github.com/babon21/excel-offer-storage/internal/offer/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGetList(t *testing.T) {
	mockOfferRepo := new(mocks.OfferRepository)
	mockOfferGateway := new(mocks.OfferGateway)
	mockOfferReader := new(mocks.OfferReader)
	mockOffer := domain.Offer{
		Name:     "offerName",
		Price:    10,
		Quantity: 20,
	}

	mockListOffer := make([]domain.Offer, 0)
	mockListOffer = append(mockListOffer, mockOffer)

	t.Run("success", func(t *testing.T) {
		mockOfferRepo.On("GetList", mock.AnythingOfType("string"), mock.AnythingOfType("string"),
			mock.AnythingOfType("string")).Return(mockListOffer, nil).Once()
		u := usecase.NewOfferUseCase(mockOfferRepo, mockOfferGateway, mockOfferReader)
		list, err := u.GetList("", "", "")
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListOffer))

		mockOfferRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockOfferRepo.On("GetList", mock.AnythingOfType("string"), mock.AnythingOfType("string"),
			mock.AnythingOfType("string")).Return(nil, errors.New("Unexpected Error")).Once()
		u := usecase.NewOfferUseCase(mockOfferRepo, mockOfferGateway, mockOfferReader)
		list, err := u.GetList("", "", "")

		assert.Error(t, err)
		assert.Nil(t, list)
		mockOfferRepo.AssertExpectations(t)
	})
}

func TestStore(t *testing.T) {
	mockOfferRepo := new(mocks.OfferRepository)
	mockOfferGateway := new(mocks.OfferGateway)
	mockOfferReader := new(mocks.OfferReader)

	mockOffer := domain.Offer{
		Name:     "offerName",
		Price:    10,
		Quantity: 20,
	}

	mockListOffer := make([]domain.Offer, 0)
	mockListOffer = append(mockListOffer, mockOffer)

	t.Run("success", func(t *testing.T) {
		mockOfferGateway.On("DownloadOffers", mock.AnythingOfType("string")).
			Return("filename", nil).Once()
		mockOfferGateway.On("DeleteOffers", mock.AnythingOfType("string")).Once()
		mockOfferReader.On("ReadAll", mock.AnythingOfType("string")).Return(mockListOffer, uint32(0), nil).Once()
		mockOfferRepo.On("GetListBySellerId", mock.AnythingOfType("string")).Return([]domain.Offer{}, nil).Once()
		mockOfferRepo.On("SaveList", mock.Anything).Return(nil).Once()
		mockOfferRepo.On("DeleteList", mock.AnythingOfType("string"), mock.Anything).Return(nil).Once()
		mockOfferRepo.On("UpdateList", mock.Anything).Return(nil).Once()
		u := usecase.NewOfferUseCase(mockOfferRepo, mockOfferGateway, mockOfferReader)
		stat, err := u.Store("", "")

		assert.NotNil(t, stat)
		assert.NoError(t, err)
		//assert.Equal(t, mockRoom.Description, tempMockRoom.Description)
		mockOfferRepo.AssertExpectations(t)
	})
}
