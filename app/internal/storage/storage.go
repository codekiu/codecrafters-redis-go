package storage

import "time"

type Storage interface {
	Get(key string) (string, bool)

	Set(key string, value string)

	SetWithExpiry(key string, value string, expiry time.Duration)

	Delete(key string)
}
