package domain

type Offer struct {
	SellerId  string `json:"seller_id" db:"seller_id"`
	OfferId   string `json:"offer_id" db:"offer_id"`
	Name      string `json:"name"`
	Price     int32  `json:"price"`
	Quantity  int32  `json:"quantity"`
	Available bool   `json:"-"`
}
