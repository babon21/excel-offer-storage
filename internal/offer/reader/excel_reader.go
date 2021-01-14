package reader

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/babon21/excel-offer-storage/internal/offer/domain"
	"strconv"
)

//type excelOfferReader struct { }

func ReadAll(filePath string) ([]domain.Offer, uint32, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if len(rows) == 1 {
		return nil, 0, nil
	}

	offers := make([]domain.Offer, 0, len(rows))
	var errRowCount uint32 = 0

OUTER:
	for i := 1; i < len(rows); i++ {
		offer := domain.Offer{}
		for j, colCell := range rows[i] {
			if err := tryFillOffer(j, colCell, &offer); err != nil {
				errRowCount++
				continue OUTER
			}
		}

		if offerIsEmpty(&offer) {
			break
		}

		if !offerIsValid(&offer) {
			errRowCount++
		}

		offers = append(offers, offer)
	}

	return offers, errRowCount, err
}

func offerIsValid(offer *domain.Offer) bool {
	// TODO need implement
	return true
}

func offerIsEmpty(offer *domain.Offer) bool {
	return !(offer.OfferId != "" || offer.Name != "")
}

func tryFillOffer(num int, value string, offer *domain.Offer) error {
	// TODO do full validation
	switch num {
	case 0:
		offerId, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}

		offer.OfferId = strconv.Itoa(int(offerId))
	case 1:
		offer.Name = value
	case 2:
		price, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}

		offer.Price = int32(price)
		if offer.Price < 0 {
			return err
		}
	case 3:
		quantity, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}

		offer.Quantity = int32(quantity)
		if offer.Quantity < 0 {
			return err
		}
	case 4:
		result, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		offer.Available = result
	default:
		fmt.Println("fillOffer error!")
	}
	return nil
}

//func fillOffer(num int, value string, offer *domain.Offer) {
//	// TODO maybe in case use type with iota
//	switch num {
//	case 0:
//		offer.OfferId = value
//	case 1:
//		offer.Name = value
//	case 2:
//		strconv.Atoi(value)
//		offer.Price, _ = strconv.Atoi(value)
//	case 3:
//		offer.Quantity = value
//	case 4:
//		offer.Available = value
//	default:
//		fmt.Println("fillOffer error!")
//	}
//}

//func NewExcelOfferReader() usecase.OfferReader {
//	return &excelOfferReader{}
//}
