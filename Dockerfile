# Builder
FROM golang:1.15-alpine3.12 as builder

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /app

COPY . .

RUN make offer

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app

WORKDIR /app

EXPOSE 9090

COPY --from=builder /app/offer /app

CMD /app/offer
UPDATE offer SET name = new.name, price = new.price, quantity = new.quantity FROM (VALUES ($1, $2, $3, $4, $5),($6, $7, $8, $9, $10),($11, $12, $13, $14, $15),($16, $17, $18, $19, $20),($21, $22, $23, $24, $25),($26, $27, $28, $29, $30),($31, $32, $33, $34, $35),($36, $37, $38, $39, $40),($41, $42, $43, $44, $45),($46, $47, $48, $49, $50)) AS new (seller_id, offer_id, name, price, quantity) WHERE offer.seller_id = new.seller_id AND offer.offer_id = new.offer_id