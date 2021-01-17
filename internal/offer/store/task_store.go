package store

import (
	"fmt"
	asyncUsecase "github.com/babon21/excel-offer-storage/internal/offer/usecase/async"
	"github.com/gomodule/redigo/redis"
)

type Store struct {
	pool *redis.Pool
}

func (s Store) Get(id int64) (string, error) {
	conn := s.pool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("GET", id))
	return value, err
}

func (s Store) Set(id int64, value string) error {
	conn := s.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", id, value)
	return err
}

func (s Store) GetNewId(idFieldName string) (int64, error) {
	conn := s.pool.Get()
	defer conn.Close()

	id, err := redis.Int64(conn.Do("INCR", idFieldName))

	return id, err
}

func NewRedisStore(host string, port string) (asyncUsecase.Store, error) {
	pool := &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		},
	}

	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil {
		return nil, err
	}

	return &Store{pool: pool}, nil
}
