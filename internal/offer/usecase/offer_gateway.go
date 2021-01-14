package usecase

type OfferGateway interface {
	DownloadOffers(url string) (string, error)
	DeleteOffers(filename string)
}
