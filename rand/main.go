package rand

import (
	crand "crypto/rand"
	"encoding/binary"
	"log"
	mrand "math/rand"
)

var rnd *mrand.Rand

func init() {
	// nolint // This is explicitly OK as we hand math/rand a strong random source
	rnd = mrand.New(cryptoSource{})
}

func Intn(i int) int {
	return rnd.Intn(i)
}

func Shuffle(n int, swap func(int, int)) {
	rnd.Shuffle(n, swap)
}

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
}
