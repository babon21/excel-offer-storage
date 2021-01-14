-- create table offer

DROP TABLE IF EXISTS offer;

CREATE TABLE offer (
    seller_id INTEGER,
    offer_id INTEGER,
    name TEXT NOT NULL,
    price INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    PRIMARY KEY (seller_id, offer_id)
);

CREATE INDEX offer_offer_id_index ON offer (offer_id);
CREATE INDEX offer_seller_id_index ON offer (seller_id);
