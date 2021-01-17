package http

import (
	"github.com/babon21/excel-offer-storage/internal/offer/usecase"
	asyncUsecase "github.com/babon21/excel-offer-storage/internal/offer/usecase/async"
	"github.com/babon21/excel-offer-storage/pkg/delivery/http/api"
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// OfferHandler  represent the httphandler for offer
type OfferHandler struct {
	OfferUsecase      usecase.OfferUseCase
	AsyncOfferUsecase asyncUsecase.AsyncOfferUseCase
}

// NewOfferHandler will initialize the bookings/ resources endpoint
func NewOfferHandler(e *echo.Echo, us usecase.OfferUseCase, asyncUs asyncUsecase.AsyncOfferUseCase) {
	handler := &OfferHandler{
		OfferUsecase:      us,
		AsyncOfferUsecase: asyncUs,
	}
	e.GET("/offers", handler.GetList)
	e.POST("/offers", handler.DownloadOffers)
	e.POST("/offers/async", handler.AsyncDownloadOffers)
	e.GET("/tasks/:id", handler.GetTask)
}

// GetList will fetch the booking based on given params
func (a *OfferHandler) GetList(c echo.Context) error {
	var request api.GetListRequest
	err := c.Bind(&request)
	if err != nil {
		return c.JSONPretty(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()}, "  ")
	}

	offers, err := a.OfferUsecase.GetList(request.SellerId, request.OfferId, request.Name)
	if err != nil {
		return c.JSONPretty(getStatusCode(err), ResponseError{Message: err.Error()}, "  ")
	}

	response := api.GetListResponse{Offers: offers}
	return c.JSONPretty(http.StatusOK, response, "  ")
}

// DownloadOffers will store the room by given request body
func (a *OfferHandler) DownloadOffers(c echo.Context) (err error) {
	var request api.DownloadOffersRequest
	err = c.Bind(&request)
	if err != nil {
		return c.JSONPretty(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()}, "  ")
	}

	// TODO url validation
	//if err = isRequestValid(&request); err != nil {
	//	return c.JSONPretty(http.StatusBadRequest, ResponseError{Message: err.Error()}, "  ")
	//}

	statistic, err := a.OfferUsecase.Store(request.SellerId, request.Url)
	if err != nil {
		return c.JSONPretty(getStatusCode(err), ResponseError{Message: err.Error()}, "  ")
	}

	response := api.DownloadOffersResponse{Statistic: statistic}

	return c.JSONPretty(http.StatusOK, response, "  ")
}

func (a *OfferHandler) AsyncDownloadOffers(c echo.Context) error {
	var request api.DownloadOffersRequest
	err := c.Bind(&request)
	if err != nil {
		return c.JSONPretty(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()}, "  ")
	}

	taskId, err := a.AsyncOfferUsecase.Store(request.SellerId, request.Url)

	response := api.AsyncDownloadOffersResponse{TaskId: taskId}
	return c.JSONPretty(http.StatusOK, response, "  ")
}

func (a *OfferHandler) startDownloadOffersTask(ch chan string, sellerId string, url string) {

	// TODO get new store ID
	// TODO write taskId to ch channel
	statistic, err := a.OfferUsecase.Store(sellerId, url)
	if err != nil {
		// TODO put error to status store
	}
	_ = statistic
	// TODO put statistic and status store to store
}

func (a *OfferHandler) GetTask(c echo.Context) error {
	taskId := c.Param("id")
	strTaskId, err := strconv.ParseInt(taskId, 10, 32)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, nil, "  ")
	}
	task := a.AsyncOfferUsecase.GetTask(strTaskId)
	return c.JSONPretty(http.StatusOK, task, "  ")
}

//func isRequestValid(request *api.AddBookingRequest) error {
//	dateLayout := "2006-01-02"
//
//	_, err := time.Parse(dateLayout, request.StartDate)
//	if err != nil {
//		return errors.New("start_date field isn't valid, must be in year-month-day format")
//	}
//
//	_, err = time.Parse(dateLayout, request.ExpirationDate)
//	if err != nil {
//		return errors.New("expiration_date field isn't valid, must be in year-month-day format")
//	}
//
//	return nil
//}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	log.Error().Msg(err.Error())
	switch err {
	//case domain.ErrInternalServerError:
	//	return http.StatusInternalServerError
	//case domain.ErrNotFound:
	//	return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
