package shorten

import (
	"github.com/sqids/sqids-go"
	"math/rand"
	"time"
)

func GenerateShortKey() string {
	number := rand.New(rand.NewSource(time.Now().UnixNano())).Uint64()

	hash, _ := sqids.New()
	id, _ := hash.Encode([]uint64{number})
	return id
}
