package shorten

import (
	"github.com/sqids/sqids-go"
)

func Shorten(number uint32) string {
	hash, _ := sqids.New()
	id, _ := hash.Encode([]uint64{uint64(number)})
	return id
}
