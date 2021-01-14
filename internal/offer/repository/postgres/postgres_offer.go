package postgres

import (
	"fmt"
	"github.com/babon21/excel-offer-storage/internal/offer/domain"
	"github.com/babon21/excel-offer-storage/internal/offer/usecase"
	"github.com/jmoiron/sqlx"
	"strings"
)

type postgresOfferRepository struct {
	Conn *sqlx.DB
}

// NewPostgresOfferRepository will create an object that represent the OfferRepository interface
func NewPostgresOfferRepository(conn *sqlx.DB) usecase.OfferRepository {
	return &postgresOfferRepository{conn}
}

func (repo *postgresOfferRepository) DeleteList(sellerId string, offerIds []string) error {
	if len(offerIds) == 0 {
		return nil
	}
	return repo.bulkDelete(sellerId, offerIds)
}

func (repo *postgresOfferRepository) bulkDelete(sellerId string, offerIds []string) error {
	deleteQuery := fmt.Sprintf("DELETE FROM offer WHERE seller_id = %s AND offer_id IN (%s)", sellerId, strings.Join(offerIds, ","))
	_, err := repo.Conn.Exec(deleteQuery)
	return err
}

func (repo *postgresOfferRepository) UpdateList(offers []domain.Offer) error {
	if len(offers) == 0 {
		return nil
	}
	return repo.bulkUpdate(offers)
}

func (repo *postgresOfferRepository) bulkUpdate(offers []domain.Offer) error {
	query := createBulkUpdateQuery(offers)
	_, err := repo.Conn.Exec(query)
	return err
}

func createBulkUpdateQuery(offers []domain.Offer) string {
	valueStrings := make([]string, 0, len(offers))
	for _, offer := range offers {
		valueStrings = append(valueStrings, fmt.Sprintf("(%s, %s, '%s', %d, %d)", offer.SellerId, offer.OfferId, offer.Name, offer.Price, offer.Quantity))
	}

	return fmt.Sprintf("UPDATE offer SET name = new.name, price = new.price, quantity = new.quantity FROM (VALUES %s) AS new (seller_id, offer_id, name, price, quantity) WHERE offer.seller_id = new.seller_id AND offer.offer_id = new.offer_id",
		strings.Join(valueStrings, ","))
}

func (repo *postgresOfferRepository) GetListBySellerId(sellerId string) ([]domain.Offer, error) {
	getListQuery := "SELECT * FROM offer WHERE seller_id = $1"
	return repo.getList(getListQuery, sellerId)
}

func (repo *postgresOfferRepository) GetList(sellerId string, offerId string, offerName string) ([]domain.Offer, error) {
	getListQuery := createGetListWithOptionalQuery(sellerId, offerId, offerName)
	return repo.getList(getListQuery)
}

func createGetListWithOptionalQuery(sellerId string, offerId string, offerName string) string {
	valueStrings := make([]string, 0, 3)

	if sellerId != "" {
		valueStrings = append(valueStrings, "seller_id = "+sellerId)
	}

	if offerId != "" {
		valueStrings = append(valueStrings, "offer_id = "+offerId)
	}

	if offerName != "" {
		valueStrings = append(valueStrings, fmt.Sprintf("name LIKE '%%%s%%'", offerName))
	}

	condition := ""
	if len(valueStrings) != 0 {
		condition = strings.Join(valueStrings, " AND ")
	}

	if condition == "" {
		return "SELECT * FROM offer"
	}

	return fmt.Sprintf("SELECT * FROM offer WHERE %s", condition)
}

func (repo *postgresOfferRepository) getList(query string, args ...interface{}) ([]domain.Offer, error) {
	offers := make([]domain.Offer, 0, 1)
	err := repo.Conn.Select(&offers, query, args...)
	return offers, err
}

func (repo *postgresOfferRepository) SaveList(offers []domain.Offer) error {
	if len(offers) == 0 {
		return nil
	}
	return repo.bulkInsert(offers)
}

func (repo *postgresOfferRepository) bulkInsert(offers []domain.Offer) error {
	valueStrings := make([]string, 0, len(offers))
	valueArgs := make([]interface{}, 0, len(offers)*5)
	i := 1
	for _, offer := range offers {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i, i+1, i+2, i+3, i+4))
		valueArgs = append(valueArgs, offer.SellerId)
		valueArgs = append(valueArgs, offer.OfferId)
		valueArgs = append(valueArgs, offer.Name)
		valueArgs = append(valueArgs, offer.Price)
		valueArgs = append(valueArgs, offer.Quantity)
		i += 5
	}
	stmt := fmt.Sprintf("INSERT INTO offer (seller_id, offer_id, name, price, quantity) VALUES %s",
		strings.Join(valueStrings, ","))
	_, err := repo.Conn.Exec(stmt, valueArgs...)
	return err
}
