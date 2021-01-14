package domain

type Offer struct {
	SellerId  string `db:"seller_id"`
	OfferId   string `db:"offer_id"`
	Name      string
	Price     int32
	Quantity  int32
	Available bool
}
