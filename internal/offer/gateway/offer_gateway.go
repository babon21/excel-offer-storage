package gateway

import (
	"github.com/babon21/excel-offer-storage/internal/offer/usecase"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type offerGateway struct {
	destination string
}

func (gateway *offerGateway) DeleteOffers(filename string) {
	_ = os.Remove(filename)
}

func (gateway *offerGateway) DownloadOffers(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	file, err := ioutil.TempFile(gateway.destination, "*.xlsx")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, resp.Body)
	return file.Name(), err
}

func NewOfferGateway(destination string) usecase.OfferGateway {
	return &offerGateway{destination: destination}
}
