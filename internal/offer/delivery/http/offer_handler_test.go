package http_test

import (
	"encoding/json"
	"errors"
	offerHttp "github.com/babon21/excel-offer-storage/internal/offer/delivery/http"
	"github.com/babon21/excel-offer-storage/internal/offer/domain"
	"github.com/babon21/excel-offer-storage/internal/offer/domain/mocks"
	"github.com/babon21/excel-offer-storage/internal/offer/usecase"
	"github.com/babon21/excel-offer-storage/pkg/delivery/http/api"
	"github.com/bxcodec/faker/v3"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetList(t *testing.T) {
	var mockOffer domain.Offer
	err := faker.FakeData(&mockOffer)
	assert.NoError(t, err)
	mockUCase := new(mocks.OfferUseCase)
	mockListOffer := make([]domain.Offer, 0)
	mockListOffer = append(mockListOffer, mockOffer)

	mockUCase.On("GetList", mock.AnythingOfType("string"), mock.AnythingOfType("string"),
		mock.AnythingOfType("string")).Return(mockListOffer, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/offers", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := offerHttp.OfferHandler{
		OfferUsecase: mockUCase,
	}
	err = handler.GetList(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetListError(t *testing.T) {
	mockUCase := new(mocks.OfferUseCase)

	mockUCase.On("GetList", mock.AnythingOfType("string"), mock.AnythingOfType("string"),
		mock.AnythingOfType("string")).Return(nil, errors.New("Unexpected Error"))

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/booking?room=1", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := offerHttp.OfferHandler{
		OfferUsecase: mockUCase,
	}
	err = handler.GetList(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDownloadOffers(t *testing.T) {
	request := api.DownloadOffersRequest{
		SellerId: "1",
		Url:      "some_url",
	}

	mockUCase := new(mocks.OfferUseCase)

	j, err := json.Marshal(request)
	assert.NoError(t, err)

	mockUCase.On("Store", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(usecase.Statistic{}, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/offers", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := offerHttp.OfferHandler{
		OfferUsecase: mockUCase,
	}
	err = handler.DownloadOffers(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}
